## ScopeLens上云记录: 构建微服务(2)

![Photo by [Tolga Ulkan](https://unsplash.com/@tolga__?utm_source=medium&utm_medium=referral) on [Unsplash](https://unsplash.com?utm_source=medium&utm_medium=referral)](https://cdn-images-1.medium.com/max/10526/0*tEW0xKnBTWVw_45f)

## Preface

本节的主要内容是将组成ScopeLens的[前端、后端](https://github.com/txfs19260817/ScopeLens)和[深度学习](https://github.com/txfs19260817/Pokemon-sprites-classifier)这三个服务部署至Kubernetes环境的笔记。

我不认为ScopeLens是一个适合作为练习Kubernetes微服务部署的项目，原因包括它涉及到使用外部服务（如数据库MongoDB）等，因此本篇记录的内容最多仅达到参考价值。如果读者想练习微服务在Kubernetes上的部署，这里推荐Istio官方文档的Example: [Learn Microservices using Kubernetes and Istio](https://istio.io/latest/docs/examples/microservices-istio/)，一直到”Add a new version of reviews”之前的部分都是不涉及Istio本身的微服务搭建教程，非常有助于入门学习。

## General Procedures

将旧项目改造为适合Kubernetes环境的“现代”项目，通用的大致的流程是：

 1. 编写Dockerfile。值得一提的是，往往我们希望通过Dockerfile构建出的Image越小越好，因此写好Dockerfile也是一门技术；

 2. 构建Docker image并推送至Docker Hub（或者是其它的Docker Registry）；

 3. 创建基本的Kubernetes资源：Deployment和Service；

 4. 为应用程序添加ConfigMap；

一般编写资源配置文件的顺序是先Deployment后Service，有必要的话再写ConfigMap，有时候可能还会涉及到[RBAC](https://kubernetes.io/zh/docs/reference/access-authn-authz/rbac/)的一系列资源，不过我这里没有用到。

关于应用程序的配置问题，一个比较省事的做法往往是在编写Dockerfile的时候把配置文件也打包进去，或者是存入环境变量（在Dockerfile或Kubernetes Deployment中配置环境变量），这样做的弊端是每次对配置文件做改动后还需要重新进行构建、推送和更新Deployment来拉取镜像（tip: 命令是kubectl restart rollout deployment/<DEPLOYMENT-NAME>）的繁琐过程。因此使用ConfigMap将配置文件挂载到volume下是个不错的方式。

## Commands

本节介绍经常会用到的命令

### Docker

**构建Docker image的命令**（Docker Hub作为镜像仓库时）

    docker build -t <username>/<image>:<version> .

* -t: 镜像的tag；

* username: Docker Hub上注册的用户名；

* image: 镜像名字，必须是小写；

* version: 版本号，也可以去掉:<version>这段（冒号也去掉），这样版本号就默认为latest。

* .: 上下文，一般Dockerfile放在项目的根目录，执行命令的目录也在这里。

Trap: 如果输出提示拉取基础镜像时发生403错误，解决方法如下：

 1. 设置Docker镜像；

 2. 手动用docker pull命令把使用的到的基础镜像拉取下来。

比如在Dockerfile里使用了FROM alpine:latest指令，那么需要在运行构建命令之前执行命令docker pull alpine:latest将该镜像预先拉取至本地。

**推送镜像的命令**

    docker push <username>/<image>:<version>

如果推送成功，可以去[Docker Hub](https://hub.docker.com/)检查一下。

### Kubernetes

**创建YAML配置的资源的命令**

    kubectl apply -f <filename/path>

用YAML文件配置资源是推荐的方式。在创建全新的资源时，把apply换成create也可以，二者的区别[见此](https://stackoverflow.com/questions/47369351/kubectl-apply-vs-kubectl-create)。如果想删除所创建的资源，则把apply换成delete。

顺带一提，初学者可能会踩的坑是直接根据名称删除Pod，其实这样是无法清理资源的，因为删除后会被Deployment发现存活的Pod比desired数目要少，从而新启一个Pod。正确做法是删除Deployment。

**查看Pod状态的命令**

    kubectl get pods -n <namespace>

如果不指定namespace，则返回所有存在于default命名空间下的Pods，对于其他受命名空间约束的资源类型也是如此。如果要查看所有命名空间下的Pod情况，需使用-A标志。

正常情况下，每个Pod的状态是Running的话则表示正常。否则需要进一步使用describe和logs命令排查问题，这里就不展开篇幅了，请参阅官方文档。

## Go Build!

本节正式开始记录微服务项目的准备过程，每个项目都会被介绍到。为了避免麻烦，此处配置的所有Kubernetes资源都没有指定namespace，即它们都将被放置在default命名空间下，其实这样做不是很好。

### ScopeLens-Server

**Dockerfile**

ScopeLens的后端是Go语言编写的，Go语言的项目编译后是一个二进制文件，因此我们编写Dockerfile的策略是先用Go环境编译项目，再将得到的二进制文件送入最小的镜像里去运行。

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=dockerfile&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=%2523%2520Build%2520stage%250AFROM%2520golang%253A1.16-alpine3.13%2520AS%2520builder%250AENV%2520GOPROXY%2520https%253A%252F%252Fgoproxy.cn%252Cdirect%250AWORKDIR%2520%252Fapp%250ACOPY%2520.%2520.%250ARUN%2520go%2520build%2520-v%2520.%250A%250A%2523%2520Run%2520stage%2520%28set%2520port%2520%2560-p%2520%2524PORT%253A%2524PORT%2560%29%250AFROM%2520alpine%253A3.13%250AWORKDIR%2520%2524GOPATH%252Fsrc%252Fgithub.com%252Ftxfs19260817%252Fscopelens%252Fserver%250ACOPY%2520--from%253Dbuilder%2520%252Fapp%252Fserver%2520.%250ARUN%2520mkdir%2520config%2520log%2520assets%250ACOPY%2520assets%2520.%252Fassets%250AARG%2520PORT%253D8888%250AEXPOSE%2520%2524PORT%250A%250A%2523%2520run%250AENTRYPOINT%2520%255B%2522.%252Fserver%2522%255D)](https://cdn-images-1.medium.com/max/2000/1*n-BrDxk96dHZTcKm0MjyPg.png)

在Build stage，在Go语言的环境下构建项目，得到二进制文件。然后另起一个Run stage，把二进制文件从上一个stage复制过来，当然也要记得把所需的静态资源从本地也复制过来。最后设置启动程序的ENTRYPOINT。这样得到的Image会相当小巧，就这个而言是19.06MB，而且还算上了静态资源。

**Deployment YAML**

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=kind%253A%2520Deployment%250AapiVersion%253A%2520apps%252Fv1%250Ametadata%253A%250A%2520%2520name%253A%2520scopelens-server-deployment%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520server%250A%2520%2520%2520%2520version%253A%2520v1%250Aspec%253A%250A%2520%2520replicas%253A%25201%250A%2520%2520selector%253A%250A%2520%2520%2520%2520matchLabels%253A%250A%2520%2520%2520%2520%2520%2520app%253A%2520server%250A%2520%2520%2520%2520%2520%2520version%253A%2520v1%250A%2520%2520template%253A%250A%2520%2520%2520%2520metadata%253A%250A%2520%2520%2520%2520%2520%2520labels%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520app%253A%2520server%250A%2520%2520%2520%2520%2520%2520%2520%2520version%253A%2520v1%250A%2520%2520%2520%2520spec%253A%250A%2520%2520%2520%2520%2520%2520containers%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520scopelens-server%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520image%253A%2520%253CIMAGE%253E%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520ports%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520http%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520containerPort%253A%25208888%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520resources%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520requests%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520memory%253A%2520512Mi%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520cpu%253A%2520512m%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520limits%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520memory%253A%2520768Mi%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520cpu%253A%25201024m%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520env%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520GIN_MODE%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520value%253A%2520release%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520volumeMounts%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520server-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520mountPath%253A%2520%2522app%252Fconfig%2522%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520readOnly%253A%2520true%250A%2520%2520%2520%2520%2520%2520volumes%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520server-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520configMap%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520name%253A%2520scopelens-server-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520items%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520key%253A%2520%2522config.ini%2522%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520path%253A%2520%2522config.ini%2522)](https://cdn-images-1.medium.com/max/2000/1*WvcVxc07_4nN3Z5INqwZeQ.png)

这里梳理一下[Deployment](https://kubernetes.io/zh/docs/concepts/workloads/controllers/deployment/)甚至其他资源必备的几处配置，更详细的解说还请参阅官方文档：

* metadata: 指定好name和labels，前者是必填项（我这里还多余加上了-deployment，其实没必要，下面的资源也是），后者是自由发挥的地方，key和value合理设置好即可，便于对资源做批量选择时使用，例如查看日志。此外metadata下还支持annotations，这里没有使用到，请参阅官方文档；

* replicas: Pod数目，多启动几个副本以保证整体服务的高可用性，每个副本会被均衡地调度到每个node上；

* spec.containers.ports: 这是个数组配置，指定容器打开的所有端口，每个端口要赋予相应的名字；

* resources: Pod的资源分配。顾名思义，requests是要求的资源大小，指定后Pod会直接吃掉node上的这部分资源；limits是这个Pod最大可用的资源上限。例如Pod试图申请超出上限的内存，可能会因为OOM而被杀掉，进而重启。所以这两个地方要进行合理配置，不过怎么合理配置我也不知道。还有一点需要注意的是CPU和内存资源的单位，大小写敏感，如果没写对的话，申请的资源可能是超乎想象的，不过也不会出什么大事，只是会因为没有符合条件的node而无法被调度导致pending；

* env: 添加环境变量，这里设置了项目使用到的[gin-gonic](https://github.com/gin-gonic/gin)框架需要的环境变量，指示环境为release；

* volumeMounts和volumes: 挂载额外文件或目录的配置，这里挂载的是配置文件config.ini，二者的name要一致；volumeMounts.mountPath是挂载目录，要配置好和入口文件的相对目录关系；volumes.configMap.name是ConfigMap资源的名字，需要正确对应。更多内容见下文的ConfigMap部分。

**Service YAML**

[Service](https://kubernetes.io/zh/docs/concepts/services-networking/service/)的配置比较简单，每个服务都大同小异。

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=kind%253A%2520Service%250AapiVersion%253A%2520v1%250Ametadata%253A%250A%2520%2520name%253A%2520scopelens-server-service%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520server%250A%2520%2520%2520%2520tier%253A%2520backend%250Aspec%253A%250A%2520%2520ports%253A%250A%2520%2520%2520%2520-%2520name%253A%2520http%250A%2520%2520%2520%2520%2520%2520protocol%253A%2520TCP%250A%2520%2520%2520%2520%2520%2520port%253A%252080%250A%2520%2520%2520%2520%2520%2520targetPort%253A%25208888%250A%2520%2520selector%253A%250A%2520%2520%2520%2520app%253A%2520server%250A%2520%2520type%253A%2520ClusterIP)](https://cdn-images-1.medium.com/max/2000/1*ccFWjzJm1kUdpBkDDXySsQ.png)

* spec.ports.targetPort: 与Deployment的spec.containers.ports相对应；

* spec.ports.port: 服务真正暴露的端口，即访问服务时使用到的端口；

* spec.type: 默认就是ClusterIP，表示只能被集群内部以IP形式访问。如果你的集群是托管在云服务商的，或者是部署在了包括Docker Desktop内置的Kubernetes这样的环境，那么将其改为LoadBalancer的话，该服务会被分配一个External IP，这样该服务就可以通过<External-IP>:<spec.ports.port>从外部访问。

关于External IP，默认情况下，EKS是使用了Elastic Load Balancer（也有可能是别的类型的Load Balancer）分配一个域名，GKE分配的是一个IP地址，Docker Desktop的Kubernetes全在localhost上。

**ConfigMap YAML**

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=apiVersion%253A%2520v1%250Akind%253A%2520ConfigMap%250Ametadata%253A%250A%2520%2520name%253A%2520scopelens-server-config%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520server%250Adata%253A%250A%2520%2520config.ini%253A%2520%257C%2520%2523%2520Paste%2520the%2520content%2520below%250A%2520%2520%2520%2520Mode%2520%253D%2520release%250A%2520%2520%2520%2520...)](https://cdn-images-1.medium.com/max/2000/1*o6or0eOY1l7_z6Xg2CCXag.png)

这个更简单。metadata.name需要与Deployment中volumes.configMap.name一致。对于data一节，key是配置文件的文件名，后面跟上: |，换行缩进后把配置文件的内容粘贴在下面即可。

### ScopeLens-Website

**Dockerfile**

对于Vue.js项目，此处参阅了[官方部署指南](https://cli.vuejs.org/zh/guide/deployment.html#docker-nginx)，即先在Node.js环境下编译项目，再把编译好的dist目录挂载到NGINX服务下。

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=dockerfile&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=%2523%2520build%2520stage%250AFROM%2520node%253Alts-alpine%2520as%2520build-stage%250AWORKDIR%2520%252Fapp%250ACOPY%2520package*.json%2520.%252F%250ARUN%2520npm%2520install%250ACOPY%2520.%2520.%250ARUN%2520npm%2520run%2520build-k8s%250A%250A%2523%2520production%2520stage%250AFROM%2520nginx%253Astable-alpine%2520as%2520production-stage%250ARUN%2520chgrp%2520-R%2520root%2520%252Fvar%252Fcache%252Fnginx%2520%252Fvar%252Frun%2520%252Fvar%252Flog%252Fnginx%2520%2526%2526%2520%255C%250A%2520%2520%2520%2520chmod%2520-R%2520770%2520%252Fvar%252Fcache%252Fnginx%2520%252Fvar%252Frun%2520%252Fvar%252Flog%252Fnginx%250ACOPY%2520--from%253Dbuild-stage%2520%252Fapp%252Fdist%2520%252Fapp%250ACOPY%2520nginx.conf%2520%252Fetc%252Fnginx%252Fnginx.conf%250A)](https://cdn-images-1.medium.com/max/2000/1*wmX09toQ1Uoz6tZhJoM5Cg.png)

Build stage的最后一条指令是自定义的npm run build-k8s，这个其实是配置在了package.json的scripts下的alias，对应的命令是”build-k8s”: “vue-cli-service build --mode kubernetes”。模式kubernetes其实是影响了所加载的.env（也叫dotenv，是Vue项目使用的一种环境变量配置），这样配置之后Vue.js会加载.env, .env.kubernetes和.env.kubernetes.local三个环境变量文件。这里分享一下.env中关于后端API目录的设置：

    NODE_ENV=production
    
    VUE_APP_URL=/api/
    VUE_APP_ADVANCED_URL=/advapi/
    VUE_APP_STATIC_ASSET_URL=/assets/sprites

该项目使用的HTTP包是Axios，也就是说Axios在请求后端时的URL不再需要host部分。这也与后面NGINX的配置有些关系。

Production stage有个赋予权限的指令，是因为nginx在某些配置情况下会出现读取文件时权限不足的奇怪问题。

测试可以使用命令docker run -d -p 8080:80 <tag>然后打开浏览器访问一下localhost:8080。

**Deployment YAML**

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=kind%253A%2520Deployment%250AapiVersion%253A%2520apps%252Fv1%250Ametadata%253A%250A%2520%2520name%253A%2520scopelens-website-deployment%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520website%250A%2520%2520%2520%2520version%253A%2520v1%250Aspec%253A%250A%2520%2520replicas%253A%25201%250A%2520%2520selector%253A%250A%2520%2520%2520%2520matchLabels%253A%250A%2520%2520%2520%2520%2520%2520app%253A%2520website%250A%2520%2520%2520%2520%2520%2520version%253A%2520v1%250A%2520%2520template%253A%250A%2520%2520%2520%2520metadata%253A%250A%2520%2520%2520%2520%2520%2520labels%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520app%253A%2520website%250A%2520%2520%2520%2520%2520%2520%2520%2520version%253A%2520v1%250A%2520%2520%2520%2520spec%253A%250A%2520%2520%2520%2520%2520%2520containers%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520scopelens-website%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520image%253A%2520%253CIMAGE%253E%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520ports%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520http%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520containerPort%253A%252080%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520resources%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520requests%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520memory%253A%2520256Mi%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520limits%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520memory%253A%2520512Mi%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520volumeMounts%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520website-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520mountPath%253A%2520%252Fetc%252Fnginx%252Fnginx.conf%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520subPath%253A%2520nginx.conf%250A%2520%2520%2520%2520%2520%2520volumes%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520website-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520configMap%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520name%253A%2520scopelens-website-config)](https://cdn-images-1.medium.com/max/2000/1*S1rZKHqc1JUy0XrtJ5Szew.png)

Trap: 有些特别的地方在于多了个volumeMounts.subPath。/etc/nginx下还有别的配置文件，如果不加上这一条，那么这个目录下将变得只会有nginx.conf这一个文件，其它文件都抹除了。

**Service YAML**

并无特别之处。

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=kind%253A%2520Service%250AapiVersion%253A%2520v1%250Ametadata%253A%250A%2520%2520name%253A%2520scopelens-website-service%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520website%250A%2520%2520%2520%2520tier%253A%2520frontend%250Aspec%253A%250A%2520%2520ports%253A%250A%2520%2520%2520%2520-%2520name%253A%2520http%250A%2520%2520%2520%2520%2520%2520protocol%253A%2520TCP%250A%2520%2520%2520%2520%2520%2520port%253A%252080%250A%2520%2520%2520%2520%2520%2520targetPort%253A%252080%250A%2520%2520selector%253A%250A%2520%2520%2520%2520app%253A%2520website%250A%2520%2520type%253A%2520ClusterIP)](https://cdn-images-1.medium.com/max/2000/1*Dxi17kWMwjsWwWZre46FyQ.png)

**ConfigMap YAML**

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=2x&wm=false&code=apiVersion%253A%2520v1%250Akind%253A%2520ConfigMap%250Ametadata%253A%250A%2520%2520name%253A%2520scopelens-website-config%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520website%250Adata%253A%250A%2520%2520nginx.conf%253A%2520%257C%250A%2520%2520%2520%2520user%2520nginx%253B%250A%2520%2520%2520%2520worker_processes%25201%253B%250A%2520%2520%2520%2520error_log%2520%252Fvar%252Flog%252Fnginx%252Ferror.log%2520warn%253B%250A%2520%2520%2520%2520pid%2520%252Fvar%252Frun%252Fnginx.pid%253B%250A%2520%2520%2520%2520events%250A%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520worker_connections%25201024%253B%250A%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520http%250A%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520include%2520%2520%2520%2520%2520%2520%2520%252Fetc%252Fnginx%252Fmime.types%253B%250A%2520%2520%2520%2520%2520%2520default_type%2520%2520application%252Foctet-stream%253B%250A%250A%2520%2520%2520%2520%2520%2520log_format%2520main%2520%27%2524remote_addr%2520-%2520%2524remote_user%2520%255B%2524time_local%255D%2520%2522%2524request%2522%2520%27%250A%2520%2520%2520%2520%2520%2520%27%2524status%2520%2524body_bytes_sent%2520%2522%2524http_referer%2522%2520%27%250A%2520%2520%2520%2520%2520%2520%27%2522%2524http_user_agent%2522%2520%2522%2524http_x_forwarded_for%2522%27%253B%250A%250A%2520%2520%2520%2520%2520%2520access_log%2520%252Fvar%252Flog%252Fnginx%252Faccess.log%2520main%253B%250A%2520%2520%2520%2520%2520%2520sendfile%2520on%253B%250A%2520%2520%2520%2520%2520%2520keepalive_timeout%252065%253B%250A%2520%2520%2520%2520%2520%2520gzip%2520on%253B%250A%250A%2520%2520%2520%2520%2520%2520upstream%2520Scopelens%250A%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2523%2520scopelens-server-service%2520is%2520the%2520internal%2520DNS%2520name%2520used%2520by%2520the%2520backend%2520Service%2520inside%2520Kubernetes%250A%2520%2520%2520%2520%2520%2520%2520%2520server%2520scopelens-server-service%253B%250A%2520%2520%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520%2520%2520upstream%2520Classifier%250A%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520server%2520pokemon-classifier-service%253B%250A%2520%2520%2520%2520%2520%2520%257D%250A%250A%2520%2520%2520%2520%2520%2520server%250A%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520listen%252080%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520server_name%2520localhost%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2523%2520website%250A%2520%2520%2520%2520%2520%2520%2520%2520location%2520%252F%250A%2520%2520%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_http_version%25201.1%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520root%2520%252Fapp%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520index%2520index.html%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520try_files%2520%2524uri%2520%2524uri%252F%2520%252Findex.html%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520%2520%2520%2520%2520%2523%2520main%2520service%250A%2520%2520%2520%2520%2520%2520%2520%2520location%2520%252Fapi%252F%250A%2520%2520%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_http_version%25201.1%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_pass%2520http%253A%252F%252FScopelens%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520%2520%2520%2520%2520%2523%2520deep%2520learning%2520service%250A%2520%2520%2520%2520%2520%2520%2520%2520location%2520%252Fadvapi%252F%250A%2520%2520%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_http_version%25201.1%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_pass%2520http%253A%252F%252FClassifier%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520%2520%2520%2520%2520%2523%2520static%2520resources%250A%2520%2520%2520%2520%2520%2520%2520%2520location%2520%252Fassets%252F%250A%2520%2520%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_http_version%25201.1%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_pass%2520http%253A%252F%252FScopelens%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%257D%250A%2520%2520%2520%2520%2520%2520%2520%2520%2523%2520forbidden%2520paths%250A%2520%2520%2520%2520%2520%2520%2520%2520location%2520%7E%2520%255E%252F%28%255C.user.ini%257C%255C.htaccess%257C%255C.git%257C%255C.svn%257C%255C.project%257CLICENSE%257CREADME.md%29%250A%2520%2520%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_http_version%25201.1%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520return%2520404%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%257D%250A%250A%2520%2520%2520%2520%2520%2520%2520%2520error_page%2520404%2520%252F404.html%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520location%2520%253D%2520%252F404.html%250A%2520%2520%2520%2520%2520%2520%2520%2520%257B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520proxy_http_version%25201.1%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520root%2520%252Fusr%252Fshare%252Fnginx%252Fhtml%253B%250A%2520%2520%2520%2520%2520%2520%2520%2520%257D%250A%2520%2520%2520)](https://cdn-images-1.medium.com/max/4096/1*85ByOdMCpFrhi6SEhwyOfQ.png)

这里更多的是分享NGINX的配置问题。

* include /etc/nginx/mime.types;: 这其实与Deployment一节提到的subPath有些关联，如果没有include进来这个文件，或者这个文件压根就不存在了，会导致有些CSS丢失。Chrome下看不出来，但Firefox会温馨提示：“css 未载入，因为它的 MIME 类型 “text/plain” 不是 “text/css””；

* upstream: 指定了后端服务地址，实际上是配置Service时的name，这个域名是由集群里CoreDNS解析的。这部分配合location一起看；

* proxy_http_version 1.1: 每个location下都加了一行HTTP版本的配置，是因为之后我们会使用到Istio配置Service Mesh，如果没有这一条，透过Istio Envoy访问前端时会返回状态426 Upgrade Required，这是因为Envoy 默认要求使用 HTTP/1.1 或 HTTP/2，客户端如果使用如1.0这样的低版本就会发生错误。

### Pokemon-Sprites-Classifier

这个是ScopeLens用到的深度学习服务，通过Pytorch+Flask+Gunicorn实现，它可以自动从上传的队伍分享页截图中识别出6只宝可梦的种类，加快了用户上传队伍信息的效率。

**Dockerfile**

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=dockerfile&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=FROM%2520python%253A3.8.10-slim-buster%250A%250AWORKDIR%2520%252Fusr%252Fsrc%252Fapp%250A%250A%2523%2520Install%2520dependencies%250ACOPY%2520requirements.prod.txt%2520.%250ARUN%2520pip3%2520install%2520--upgrade%2520pip%250ARUN%2520pip3%2520install%2520--no-cache-dir%2520-i%2520https%253A%252F%252Fpypi.tuna.tsinghua.edu.cn%252Fsimple%2520-r%2520requirements.prod.txt%250ARUN%2520pip3%2520install%2520--no-cache-dir%2520torch%253D%253D1.8.0%252Bcpu%2520torchvision%253D%253D0.9.0%252Bcpu%2520-f%2520https%253A%252F%252Fdownload.pytorch.org%252Fwhl%252Ftorch_stable.html%250A%250A%2523%2520Copy%2520necessary%2520files%250ARUN%2520mkdir%2520dataset%2520configs%250ACOPY%2520utils%2520.%252Futils%250ACOPY%2520app.py%2520*.pth%2520.%252F%250ACOPY%2520dataset%252Flabel.csv%2520dataset%252Flabel.csv%250A%250A%2523%2520Run%2520%28-p%2520%2524PORT%253A%2524PORT%29%250AARG%2520PORT%253D14514%250AEXPOSE%2520%2524PORT%250AENTRYPOINT%2520%255B%2520%2522gunicorn%2522%252C%2520%2522app%253Aapp%2522%252C%2520%2522-c%2522%252C%2520%2522.%252Fconfigs%252Fgunicorn.conf.py%2522%255D)](https://cdn-images-1.medium.com/max/2048/1*JaBu3HB4Yar_apVbypNKVg.png)

我不清楚Python项目编写Dockerfile的最佳实践，对此有了解的读者欢迎提出建议。总之是选择了一个较小的Python镜像，装好必要的依赖和静态文件，通过gunicorn启动即可。

Trap: 对于Pytorch项目，不要使用以Alpine为基础镜像的Python镜像，因为该发行版没有glibc，导致import torch时直接报错。

**Deployment YAML**

![[Click me to the code snippet](https://carbon.now.sh/?bg=rgba%28171%2C+184%2C+195%2C+1%29&t=nord&wt=none&l=yaml&ds=true&dsyoff=20px&dsblur=68px&wc=true&wa=true&pv=56px&ph=56px&ln=false&fl=1&fm=Hack&fs=14px&lh=133%25&si=false&es=1x&wm=false&code=kind%253A%2520Deployment%250AapiVersion%253A%2520apps%252Fv1%250Ametadata%253A%250A%2520%2520name%253A%2520pokemon-classifier-deployment%250A%2520%2520labels%253A%250A%2520%2520%2520%2520app%253A%2520classifier%250A%2520%2520%2520%2520version%253A%2520v1%250Aspec%253A%250A%2520%2520replicas%253A%25201%250A%2520%2520selector%253A%250A%2520%2520%2520%2520matchLabels%253A%250A%2520%2520%2520%2520%2520%2520app%253A%2520classifier%250A%2520%2520%2520%2520%2520%2520version%253A%2520v1%250A%2520%2520template%253A%250A%2520%2520%2520%2520metadata%253A%250A%2520%2520%2520%2520%2520%2520labels%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520app%253A%2520classifier%250A%2520%2520%2520%2520%2520%2520%2520%2520version%253A%2520v1%250A%2520%2520%2520%2520spec%253A%250A%2520%2520%2520%2520%2520%2520containers%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520sprites-classifier%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520image%253A%2520%253CIMAGE%253E%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520ports%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520http%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520containerPort%253A%252014514%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520resources%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520limits%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520memory%253A%25202048Mi%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520cpu%253A%25201024m%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520requests%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520memory%253A%25201024Mi%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520cpu%253A%2520512m%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520livenessProbe%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520failureThreshold%253A%25203%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520httpGet%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520path%253A%2520%252Fadvapi%252Fpredict%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520port%253A%252014514%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520scheme%253A%2520HTTP%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520initialDelaySeconds%253A%252030%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520periodSeconds%253A%25203600%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520timeoutSeconds%253A%252030%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520volumeMounts%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520classifier-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520mountPath%253A%2520%2522%252Fusr%252Fsrc%252Fapp%252Fconfigs%2522%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520readOnly%253A%2520true%250A%2520%2520%2520%2520%2520%2520volumes%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520-%2520name%253A%2520classifier-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520configMap%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520name%253A%2520pokemon-classifier-config%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520items%253A%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520key%253A%2520%2522config.json%2522%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520path%253A%2520%2522config.json%2522%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520-%2520key%253A%2520%2522gunicorn.conf.py%2522%250A%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520%2520path%253A%2520%2522gunicorn.conf.py%2522%250A)](https://cdn-images-1.medium.com/max/2000/1*mqDIDQZ7hVlYhWOdKOvx2g.png)

没有什么特别之处，但是想引入的一个新概念是[存活探测器livenessProbe](https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)。kubelet会以指定的方法、路径和端口访问你的服务，如果在规定时间内返回200（其他2xx或许也可以？），则kubelet知道这个Pod是存活的，否则是异常的。

Service和ConfigMap也没什么特别之处就不贴了，所有这些文件在ScopeLens的repository下都可以找到。

## Test

这一节没有演示，看文字想象一下吧。

鉴于我们的环境是EKS，因此直接把前端Scopelens-website的Service的type改为LoadBalancer（既可以修改文件并重新kubectl apply，又可以使用命令kubectl edit svc scopelens-website-service进行修改）。然后，通过kubectl get svc scopelens-website-service找到External IP，打开浏览器访问试试即可。注意安全组要打开80端口（顺便可以把443也打开），EKS需要配置的安全组是“其他安全组”，在EKS网页控制台可以找到。

此外，也可以用kubectl exec命令，在一个Pod上使用curl来测试以上所有的服务。这个在[Istio示例的生产测试一节](https://istio.io/latest/docs/examples/microservices-istio/production-testing/)也有提到，他用到了一个sleep的Pod。

## Conclusion

踩了很多坑终于可以玩起来了。但是离生产环境还有几个问题等待解决：

* 域名是云服务商的Load Balancer分配的，拿不出手而且不是固定的；

* 所有服务都是HTTP协议的，没有TLS加密；

* 只将前端的Service配置为LoadBalancer时，如果想通过分配的域名直接访问后端服务的话，不可避免地要走一遍前端的NGINX，产生了耦合；而3个服务都开启LoadBalancer的话，要记住3个奇怪的域名。

其实有些问题是依靠Ingress解决的，这篇文章没有像其他Kubernetes入门教程一样配置Ingress，这是因为我打算一步到位，用一些额外的云原生应用进解决问题。这里提前揭露一下后续涉及到的内容：Istio, External DNS和Cert-manager。我们下次再见。

## Reference

 1. [Learn Microservices using Kubernetes and Istio](https://istio.io/latest/docs/examples/microservices-istio/)

 2. [K8S中挂载目录引发的血案！](https://zhuanlan.zhihu.com/p/340548846)

 3. [istio 常见问题: 返回426 状态码](https://imroc.cc/post/202105/426-status-code/)