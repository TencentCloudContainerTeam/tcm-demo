const util = require('util')
const Koa = require('koa')
const app = new Koa()
const request = require("request")
const requestPromise = util.promisify(request);

app.use( async ( ctx ) => {
  let url = ctx.request.url

  let discountOptions = {
    uri: "http://127.0.0.1:4000/discount",
    method: "GET",
    json: true
  }
  let discount = await requestPromise(discountOptions);

  let recommendOptions = {
    uri: "http://127.0.0.1:3000/recommend",
    method: "GET",
    json: true
  }
  let recommend = await requestPromise(recommendOptions);

  ctx.body = JSON.stringify(discount.body) + '\n\n' + JSON.stringify(recommend.body)
})
app.listen(8000)
