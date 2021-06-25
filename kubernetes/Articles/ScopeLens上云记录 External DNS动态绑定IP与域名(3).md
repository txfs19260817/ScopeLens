
## ScopeLens上云记录: External DNS动态绑定IP与域名(3)

![Photo by [Nadine Shaabana](https://unsplash.com/@nadineshaabana?utm_source=medium&utm_medium=referral) on [Unsplash](https://unsplash.com?utm_source=medium&utm_medium=referral)](https://cdn-images-1.medium.com/max/8000/0*NseE2g4DLpj6oCVn)

## How is ExternalDNS useful to me?

### Problem

这是[ExternalDNS](https://github.com/kubernetes-sigs/external-dns)项目[FAQ](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/faq.md)的第一个问题“它有什么用？”。上一篇文章提到，通过LoadBalancer暴露服务后，服务会被云服务商分配一个External IP，可能是IP也可能是域名，外部客户端可以凭此访问集群上的服务。但是这个External IP不美观很难记，而绝大多数情况下我们都希望用花钱买的域名（FQDN）来公开服务。这时解决方案就浮出水面了——我们在DNS解析商添加一条域名到External IP的A记录不就可以了？但是External IP的第二个特点是不固定，可能你重启一下Service/Ingress，就又会被分配一个新External IP，这时域名就访问不到服务了。

### Those times are over!

简短截说，ExternalDNS自动在DNS解析商帮你维护域名到External IP的记录。ExternalDNS会定期检测新以LoadBalancer的形式暴露的Service、新的配置了hosts的Ingress资源的创建、任何External IP发生变动等事件，一旦检测到就会更新DNS解析，让域名永远指向最新可用的External IP。ExternalDNS[支持的DNS Providers见此](https://github.com/kubernetes-sigs/external-dns#status-of-providers)。

## Setup

由于我们的集群是托管于EKS的，所以域名的解析服务这里也选择了AWS Route 53。不足之处在于它不是免费的。顺带一提，Route 53有出售域名的服务，本站的新域名slen.cc也是从这里购买的。如果你已经拥有一个域名，可以尝试将其迁移到Route 53上去托管。托管域名的单位叫托管区。

![托管在Route 53的域名](https://cdn-images-1.medium.com/max/2000/1*IvRBt4eUbmMbd9vA4Nb3zw.png)

### IAM Policy

第一步不是安装，而是先配置好一个附加了可以控制Route 53的IAM Policy的角色，External DNS的Service Account要“扮演”这个角色以可以操控Route 53来增加、修改DNS记录。

一般来说，我们赋予服务最小的权限范围是一种安全的做法，External DNS文档给出的IAM Policy如下：

    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "route53:ChangeResourceRecordSets"
          ],
          "Resource": [
            "arn:aws:route53:::hostedzone/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "route53:ListHostedZones",
            "route53:ListResourceRecordSets"
          ],
          "Resource": [
            "*"
          ]
        }
      ]
    }

你可以前往AWS IAM网页控制台，创建一份策略，把该文本粘贴进去，并取一个名字，例如官方文档将其命名为AllowExternalDNSUpdates，保存后会得到它的ARN，形如arn:aws:iam::12345678:policy/AllowExternalDNSUpdates。

接着借助eksctl来创建一个[IAM Roles for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)，实现Kubernetes的Service Account资源与IAM Role的映射。

    eksctl create iamserviceaccount --cluster=<CLUSTER-NAME> --name=external-dns --attach-policy-arn=<IAM-POLICY-ARN> --override-existing-serviceaccounts --approve

此时查看IAM控制台，会发现多了一个由eksctl创建的角色，把它的ARN先记下来，稍后会用到。

### Set up a hosted zone

继续后面的工作之前，如果没有安装[jq](https://github.com/stedolan/jq)的话请安装一下。

这时假设已经拥有一个域名并由Route 53托管区托管。如果没有的话可以执行类似下面的命令创建一个：

    aws route53 create-hosted-zone --name "external-dns-test.my-org.com." --caller-reference "external-dns-test-$(date +%s)"

注意将”external-dns-test.my-org.com.”引号内的域名换成在托管区里的那个域名，最后的.不要漏掉。

然后获取hostedzone ID：

    aws route53 list-hosted-zones-by-name --output json --dns-name "external-dns-test.my-org.com." | jq -r '.HostedZones[0].Id'

控制台会输出一个形如/hostedzone/ZEWFWZ4R16P7IB的ID，它和网页控制台上的ID（随机字符部分）应该是一样的，记录下来后面会用到。

官方文档还多了一步列出nameservers的命令，以便之后用dig命令测试，我这里选择直接打开网页测试所以没有用到，不过还是把命令列出来，测试的部分请直接参阅文档（见文末）。

    aws route53 list-resource-record-sets --output json --hosted-zone-id "/hostedzone/ZEWFWZ4R16P7IB" \
        --query "ResourceRecordSets[?Type == 'NS']" | jq -r '.[0].ResourceRecords[].Value'

### Deploy ExternalDNS

![[Click to src](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=apiVersion%253A%2520v1%250Akind%253A%2520ServiceAccount%250Ametadata%253A%250A%2520%2520name%253A%2520external-dns%250A%2520%2520%2523%2520If%2520you%27re%2520using%2520Amazon%2520EKS%2520with%2520IAM%2520Roles%2520for%2520Service%2520Accounts%252C%2520specify%2520the%2520following%2520annotation.%250A%2520%2520%2523%2520Otherwise%252C%2520you%2520may%2520safely%2520omit%2520it.%250A%2520%2520annotations%253A%250A%2520%2520%2520%2520%2523%2520Substitute%2520your%2520account%2520ID%2520and%2520IAM%2520service%2520role%2520name%2520below.%250A%2520%2520%2520%2520eks.amazonaws.com%252Frole-arn%253A%2520arn%253Aaws%253Aiam%253A%253A%253CACCOUNT-ID%253E%253Arole%252F%253CIAM-SERVICE-ROLE-NAME%253E%250A---%250AapiVersion%253A%2520rbac.authorization.k8s.io%252Fv1%250Akind%253A%2520ClusterRole%250Ametadata%253A%250A%2520%2520name%253A%2520external-dns%250Arules%253A%250A-%2520apiGroups%253A%2520%255B%2522%2522%255D%250A%2520%2520resources%253A%2520%255B%2522services%2522%252C%2522endpoints%2522%252C%2522pods%2522%255D%250A%2520%2520verbs%253A%2520%255B%2522get%2522%252C%2522watch%2522%252C%2522list%2522%255D%250A-%2520apiGroups%253A%2520%255B%2522extensions%2522%252C%2522networking.k8s.io%2522%255D%250A%2520%2520resources%253A%2520%255B%2522ingresses%2522%255D%250A%2520%2520verbs%253A%2520%255B%2522get%2522%252C%2522watch%2522%252C%2522list%2522%255D%250A-%2520apiGroups%253A%2520%255B%2522%2522%255D%250A%2520%2520resources%253A%2520%255B%2522nodes%2522%255D%250A%2520%2520verbs%253A%2520%255B%2522list%2522%252C%2522watch%2522%255D%250A---%250AapiVersion%253A%2520rbac.authorization.k8s.io%252Fv1%250Akind%253A%2520ClusterRoleBinding%250Ametadata%253A%250A%2520%2520name%253A%2520external-dns-viewer%250AroleRef%253A%250A%2520%2520apiGroup%253A%2520rbac.authorization.k8s.io%250A%2520%2520kind%253A%2520ClusterRole%250A%2520%2520name%253A%2520external-dns%250Asubjects%253A%250A-%2520kind%253A%2520ServiceAccount%250A%2520%2520name%253A%2520external-dns%250A%2520%2520namespace%253A%2520default%250A---%250AapiVersion%253A%2520apps%252Fv1%250Akind%253A%2520Deployment%250Ametadata%253A%250A%2520%2520name%253A%2520external-dns%250Aspec%253A%250A%2520%2520strategy%253A%250A%2520%2520%2520%2520type%253A%2520Recreate%250A%2520%2520selector%253A%250A%2520%2520%2520%2520matchLabels%253A%250A%2520%2520%2520%2520%2520%2520app%253A%2520external-dns%250A%2520%2520template%253A%250A%2520%2520%2520%2520metadata%253A%250A%2520%2520%2520%2520%2520%2520labels%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520app%253A%2520external-dns%250A%2520%2520%2520%2520%2520%2520%2523%2520If%2520you%27re%2520using%2520kiam%2520or%2520kube2iam%252C%2520specify%2520the%2520following%2520annotation.%250A%2520%2520%2520%2520%2520%2520%2523%2520Otherwise%252C%2520you%2520may%2520safely%2520omit%2520it.%250A%2520%2520%2520%2520%2520%2520annotations%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520iam.amazonaws.com%252Frole%253A%2520arn%253Aaws%253Aiam%253A%253A%253CACCOUNT-ID%253E%253Arole%252F%253CIAM-SERVICE-ROLE-NAME%253E%250A%2520%2520%2520%2520spec%253A%250A%2520%2520%2520%2520%2520%2520serviceAccountName%253A%2520external-dns%250A%2520%2520%2520%2520%2520%2520containers%253A%250A%2520%2520%2520%2520%2520%2520-%2520name%253A%2520external-dns%250A%2520%2520%2520%2520%2520%2520%2520%2520image%253A%2520k8s.gcr.io%252Fexternal-dns%252Fexternal-dns%253Av0.7.6%250A%2520%2520%2520%2520%2520%2520%2520%2520args%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--source%253Dservice%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--source%253Dingress%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--domain-filter%253D%253CDOMAIN-NAME%253E%2520%2523%2520will%2520make%2520ExternalDNS%2520see%2520only%2520the%2520hosted%2520zones%2520matching%2520provided%2520domain%252C%2520omit%2520to%2520process%2520all%2520available%2520hosted%2520zones%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--provider%253Daws%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--policy%253Dupsert-only%2520%2523%2520would%2520prevent%2520ExternalDNS%2520from%2520deleting%2520any%2520records%252C%2520omit%2520to%2520enable%2520full%2520synchronization%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--aws-zone-type%253Dpublic%2520%2523%2520only%2520look%2520at%2520public%2520hosted%2520zones%2520%28valid%2520values%2520are%2520public%252C%2520private%2520or%2520no%2520value%2520for%2520both%29%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--registry%253Dtxt%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520--txt-owner-id%253D%253CHOSTED-ZONE-ID%253E%250A%2520%2520%2520%2520%2520%2520securityContext%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520fsGroup%253A%252065534%2520%2523%2520For%2520ExternalDNS%2520to%2520be%2520able%2520to%2520read%2520Kubernetes%2520and%2520AWS%2520token%2520files)](https://cdn-images-1.medium.com/max/2048/1*P8J1M-dSUZpI4nXeEB258g.png)

* namespace: 注意这里所有资源直接部署至default命名空间下，如果有部署在其他namespace下的需求，可以**尝试**（我没验证过）：先去IAM控制台找到刚才创建的角色，然后依次点击信任关系→编辑信任关系，最后把system:serviceaccount:default:external-dns这样的字符串中的default部分改成你要设定的命名空间并保存；

* annotations：ServiceAccount和Deployment的annotations都要求设置eksctl所创建的角色的ARN；

* args的domain-filter：决定了External DNS可见的域名，可以填上你的域名，也可以留空；

* args的txt-owner-id：写上刚才控制台输出的hostedzone ID，包括/hostedzone/；

保存该文件并通过kubectl apply令其生效。

### Take effects

暴露Service或Ingress类型，让EKS赋予它们External IP的同时让External DNS更新DNS记录。

**Service**：把类型设置为LoadBalancer，并添加annotations：external-dns.alpha.kubernetes.io/hostname，值是你的域名，例如：

    apiVersion: v1
    kind: Service
    metadata:
      name: nginx
      annotations:
        external-dns.alpha.kubernetes.io/hostname: nginx.external-dns-test.my-org.com
    spec:
      type: LoadBalancer
      ports:
      - port: 80
        name: http
        targetPort: 80
      selector:
        app: nginx

**Ingress**：在annotations指定ingress controller，例如：

    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: foo
      annotations:
        kubernetes.io/ingress.class: "nginx" # use the one that corresponds to your ingress controller.
    spec:
      rules:
      - host: foo.bar.com
        http:
          paths:
          - backend:
              serviceName: foo
              servicePort: 80

等一会儿再看看Route 53控制台发现各有一条A记录和TXT记录被插入到托管区的记录里，A记录正是域名到LoadBalancer的External IP的映射。再用自己的域名访问试试，hopefully，你已经可以通过自己的域名访问到服务了。

## Conclusion

其实本文是照搬了External DNS的文档而已。下一节简介Istio的引入，我们将利用它提供的Ingress Controller来控制进入集群的流量。届时，External DNS还会有个更新，来监测Istio的Ingress的配置。我们下期再见。

## Reference

 1. [https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md)

 2. [https://docs.aws.amazon.com/zh_cn/Route53/latest/DeveloperGuide/getting-started.html](https://docs.aws.amazon.com/zh_cn/Route53/latest/DeveloperGuide/getting-started.html)
