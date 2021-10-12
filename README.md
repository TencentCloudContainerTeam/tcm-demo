注意：

* 如果使用 TCM 1.10.3：请使用 tag tcm-1.10.3
* 如果使用开源 istio：请使用 master 分支

##  在线商城系统

<img src="http://ww4.sinaimg.cn/large/006tNc79gy1g4dtmtlwdwj30o50vcwjy.jpg" referrerpolicy="no-referrer"/>

这是一个在线电子商城, 属于多语言微服务系统, 包括nodejs, golang, ruby, python和java, 整个系统由大概10个微服务项目组成.

整个页面分三个区域: 用户信息, 推荐商品, 折扣商品:

用户访问商城首页时, mall服务会分别访问后端的服务users, recommend和discount, 分别获取用户信息, 推荐商品列表, 和折扣商品列表. 其中users服务又会从mongodb服务读写用户信息, recommend服务会从scores服务中查询商品的综合评分, 同时调用products服务获取商品详情, discount服务也会从products商品服务获取商品详情.

------

## 多分支环境场景

<img src="http://ww4.sinaimg.cn/large/006tNc79gy1g4dtn2750sj31cx0u01hd.jpg" referrerpolicy="no-referrer"/>

- jason希望验证推荐系统recommend 新版本v2, 这个版本在「推荐商品」区域上增加了一个banner.
- 与此同时, fox正打算验证另一个的feature: 包括对discount 和 products的修改, 同时引入了新的收藏服务favorites. 其中discount v2在「折扣商品」区域上新增一个banner, 同时products服务会通过调用新的服务favorites获取商品的收藏人数, 然后返回给前端页面.

