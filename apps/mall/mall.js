const util = require('util')
const Koa = require('koa')
const app = new Koa()
const request = require("request")
const requestPromise = util.promisify(request);
const port = 7000

app.use( async ( ctx ) => {
  let url = ctx.request.url
  let headers = getForwardHeaders(ctx.request)
  console.log(`getting mall page for user ${headers.cookie}`) //todo

  let discountOptions = {
    uri: "http://discount.base.svc.cluster.local:7000/discount",
    method: "GET",
    headers: headers,
    json: true
  }
  let discount = await requestPromise(discountOptions);

  let recommendOptions = {
    uri: "http://recommend.base.svc.cluster.local:7000/recommend",
    method: "GET",
    headers: headers,
    json: true
  }
  let recommend = await requestPromise(recommendOptions);

  ctx.body = JSON.stringify(discount.body) + '\n\n' + JSON.stringify(recommend.body)
})

function getForwardHeaders(request) {
  headers = {}
  forwardHeaders = [
    "cookie",
    "x-request-id",
    "x-b3-traceid",
    "x-b3-spanid",
    "x-b3-parentspanid",
    "x-b3-sampled",
    "x-b3-flags",
    "x-ot-span-context",
  ]

  for (let i = 0, len = forwardHeaders.length; i < len; i++)  {
    let key = forwardHeaders[i]
    if (request.headers[key]) {
      headers[key] = request.headers[key]
    }
  }

  return headers
}

console.log(`Starting scores service on port ${port}`)
app.listen(port)
