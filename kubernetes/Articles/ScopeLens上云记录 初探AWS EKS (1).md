
## ScopeLens上云记录: 初探AWS EKS (1)

![Photo by [marina](https://unsplash.com/@marinajune?utm_source=medium&utm_medium=referral) on [Unsplash](https://unsplash.com?utm_source=medium&utm_medium=referral)](https://cdn-images-1.medium.com/max/10368/0*mB3Rxd04iAiGvg7u)

## Introduction

大家好，由此是网页应用[ScopeLens](https://github.com/txfs19260817/ScopeLens)部署至Kubernetes环境的系列文章，内容大致包括微服务的部署、服务外部暴露、日志系统构建等。常用的参考链接会在每篇文章的文末给出。有任何疑问或其它的讨论内容欢迎直接留言或发送至[邮箱](mailto:scopelens@pm.me)与我联系。此时，你已经可以通过访问[https://slen.cc/](https://slen.cc/)使用本应用。

### Motivation

* 提高自己的姿势水平；

* AWS还有一些credit不用就到期了。

## Amazon EKS

Amazon Elastic Kubernetes Service (Amazon EKS) 是一种托管的Kubernetes服务。省去了自己创建若干节点并将它们组织成Kubernetes集群（例如创建多个EC2实例并使用kubeadm建立集群）以及集群管理的麻烦。EKS的一部分优势如下：

* 自动运行和缩扩容控制平面，保证集群的高可用；

* 很容易配置节点的实例类型与数目；

* 提供一个默认的Load Balancer，很容易地将Service暴露到外部；

缺点主要就是贵，除了要为EC2实例支付费用，EKS服务本身也不是免费的，每一个EKS集群每小时需支付0.10 USD。

不过有一个节省支出的办法是启用Spot实例作为节点。这时你可能会有疑问：Spot实例如果被停机该怎么办？答案是EKS自动会在Spot实例将要被回收的前几分钟启动新的Spot实例，并将工作负载搬运至上面。如果想要了解更多，本文的最后一节给出了相关文章，本系列文章的后续部分也会讲到如何启用Spot实例作为节点。

### Trap: Pod limit on Node

本系列文章会时不时穿插上云过程中遇到的一些坑。第一个坑是关于EKS集群每个节点上可容纳的最大Pod数目，不像Google Kubernetes Engine (GKE)那样每个节点可以最多容纳110个Pods，该数目取决于EC2实例类型，每种实例的Pod容量上限取决于它的Elastic Network Interfaces (ENI)数目和每个ENI的IP数目，对应表在[Available Ip Per ENI](https://github.com/awslabs/amazon-eks-ami/blob/master/files/eni-max-pods.txt)。如果出于某种原因（比如省钱）你为ECK集群挑选了如t3.micro这种免费套餐包含的实例类型，那么需要注意，每个t3.micro实例最多只能放4个Pods，连容下kube-system本身的Pod都够呛，问题处理起来也很麻烦。所以处于性价比的考虑，尽量选择large这个级别的实例。后文我也会给出我的初步实践。

## Eksctl

顾名思义[eksctl](http://eksctl.io)是一个命令行工具，由[Weaveworks](https://www.weave.works/)开发，它可以非常方便地创建、管理以及销毁EKS集群及相关资源，实质是通过AWS CloudFormation对资源进行管理的。一个简单的eksctl create cluster命令，就可以帮你建立起一个含两个m5.large实例的EKS集群，不过通常是借助文件来配置集群的。关于eksctl的使用请多多参阅官方文档。

### Prerequisites

这里会提到开始我们的集群搭建工作之前需要的一些前置条件。

* kubectl: Kubernetes集群的管理工具。除了如何使用该命令行工具，一些关于Kubernetes的前置知识也是很有必要的，文末给出了Kubernetes基础教学的官方文档，书籍的话我比较推荐[Kubernetes权威指南](https://book.douban.com/subject/35458432/)，目前已经出到第5版，讲解版本也是较新的v1.19。kubectl的[安装指南见此](https://docs.amazonaws.cn/en_us/eks/latest/userguide/install-kubectl.html)。

* eksctl: eksctl的[安装指南见此](https://docs.amazonaws.cn/en_us/eks/latest/userguide/eksctl.html)。

* Amazon CLI: 安装Amazon CLI以启用aws命令行工具。目前只需要[安装](https://docs.amazonaws.cn/en_us/cli/latest/userguide/cli-chap-install.html)（建议使用Version 2）和运行aws configure进行[配置](https://docs.amazonaws.cn/en_us/cli/latest/userguide/cli-configure-quickstart.html#cli-configure-quickstart-config)。配置前需要在AWS IAM控制界面创建一个用户（user，还有一种类型叫角色role，后面会频繁与其打交道）。关于创建的用户权限，我这里直接创建了一个具有Administrator权限的用户，但创建用户与角色的原则是权限越小越好。创建后记得保存好用户的Access Key ID和Secret Access Key。

### Config File

将你的配置按[官方文档](https://eksctl.io/usage/schema/)写入文件，如cluster.yaml，然后运行命令eksctl create cluster -f cluster.yaml即可启动建立EKS集群的过程。

下面分享一下我的YAML配置文件：

    apiVersion: eksctl.io/v1alpha5
    kind: ClusterConfig
    
    metadata:
      name: scopelens
      region: us-west-1
      
    iam:
      withOIDC: true
    
    managedNodeGroups:
    - name: on-demand
      minSize: 1
      maxSize: 1
      desiredCapacity: 1
      volumeSize: 20
      ssh:
        allow: true
        publicKeyName: PEMKEYFILE
      instanceTypes: ["t3.medium"]
    
    - name: spot
      minSize: 2
      maxSize: 3
      desiredCapacity: 3
      volumeSize: 20
      ssh:
        allow: true
        publicKeyName: PEMKEYFILE
      instanceTypes: ["c3.large","c4.large","c5.large","c5d.large","c5n.large","c5a.large"]
      spot: true

* metadata部分给出了集群名和区域，之后在使用eksctl的过程中，集群名是flag --cluster=name经常需要的参数，而区域参数有时候也会需要，但可以不在命令行给出，eksctl会自动使用你在aws configure配置时的默认区域，这也表示如果你创建的集群所在区域与你配置aws命令行工具时指定的区域不一致，而且也未在命令行通过参数--region明确给出，则会提示找不到你的集群。

* iam.withOIDC: 将来为启用一部分功能，需要使用ServiceAccount这个Kubernetes资源，而要在EKS使用它们往往需要将它们与IAM进行绑定，使用IAM就需要启用OIDC认证。关于这部分的说明详见文末链接。

* managedNodeGroups: 节点按节点组划分与管理。这里使用的是managed node group，同样也有unmanaged类型，后者创建的节点不会显示在EKS控制台，有关二者的更多区别可以见[Reddit的讨论](https://www.reddit.com/r/aws/comments/ei87ch/im_confused_about_amazon_eks_nodegroups_vs/)。

* node group下创建了两个节点组，一个是按需实例组，另一个是Spot实例组，也就是说在节点组下配置spot: true，eksctl可以自动为你创建和管理Spot实例的节点组，也就是可以省钱的意思。

* instanceTypes: 该组可挑选实例的列表，eksctl创建节点时会从这里挑选实例类型，优先度似乎是哪个在AWS区域上剩余的多就用哪个。注意，有些实例类型在某些区域（Region）是不存在的，一定要先在[定价页面](https://aws.amazon.com/cn/ec2/pricing/on-demand/)查阅清楚再配置，否则会出错；还有就是Spot实例有时在某些可用区（Availability Zones）是不可用的，需要添加availabilityZones:参数指定可用区，否则创建会报错，至于怎么查某个实例类型的Spot实例在哪个可用区可用，可以使用[Dry Run功能](https://eksctl.io/usage/dry-run/)先试试，或者去EC2创建实例导航页面去查看。

* ssh: 是否启用SSH以允许你通过SSH和密钥连上节点实例，这个密钥必须在该区域内已经存在，关于密钥的创建与导入请参阅[Amazon EC2 密钥对和 Linux 实例](https://docs.aws.amazon.com/zh_cn/AWSEC2/latest/UserGuide/ec2-key-pairs.html)文档。

* 每组的其它配置包括节点的数目范围、挂载的EBS卷大小等。

成功创建之后，你的~/.kube/config会被更新你的集群配置，这时用kubectl get nodes方可查看所有节点信息。

### IAM Identity Mapping

[IAM Identity Mapping功能](https://eksctl.io/usage/iam-identity-mappings/)可以赋予用户user和角色role权限来管理EKS集群。通过命令eksctl create iamidentitymapping --cluster <集群名> --arn arn:aws:iam::<数字ID>:user/<用户名> --group system:masters --username <用户名> --region <区域>为你刚才创建的用户赋予权限，以后使用AWS网页控制台的时候都建议使用这个用户而非你的根账户。

## Conclusion

漂亮，这时你已经拥有一个AWS托管的Kubernetes集群了。如果你想清理资源，只需要运行命令eksctl delete cluster -f cluster.yaml即可。

下一篇文章预定讲解ScopeLens的3个微服务的部署过程。

## Reference

 1. [https://github.com/txfs19260817/ScopeLens](https://github.com/txfs19260817/ScopeLens)

 2. [Building for Cost optimization and Resilience for EKS with Spot Instances](https://aws.amazon.com/cn/blogs/compute/cost-optimization-and-resilience-eks-with-spot-instances/)

 3. [https://eksctl.io/](https://eksctl.io/)

 4. [Getting started with Amazon EKS — eksctl](https://docs.amazonaws.cn/en_us/eks/latest/userguide/getting-started-eksctl.html)

 5. [Kubernetes Basics](https://kubernetes.io/zh/docs/tutorials/kubernetes-basics/)

 6. [Amazon EKS Workshop](https://www.eksworkshop.com/)

 7. [Create an IAM OIDC provider for your cluster](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html)