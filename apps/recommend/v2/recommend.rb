require 'sinatra'
require 'net/http'

set :bind, '0.0.0.0'
set :port, 7000

get '/recommend' do
  content_type :json

  headers = getForwardHeaders(request.env)
  puts "getting recommend products for user #{headers["cookie"]}"

  ids = getRecommendIDsByUser() # todo pass user
  products = getProducts(ids, headers)
  products.each do |p|
    score = getScores([p["id"]], headers)
    p["score"] = score[p["id"].to_s]
  end
  recommend = {
    products: products,
    banner: "https://github.com/TencentCloudContainerTeam/tcm-demo/blob/master/assets/istio.png?raw=true"
  }

  recommend.to_json
end

def getRecommendIDsByUser()
  mockDB = (7..15).to_a
  mockDB.sample(4)
end

def getProducts(ids, headers)
  uri = URI("http://products.base.svc.cluster.local:7000/products?ids=#{ids.join(',')}")
  req = Net::HTTP::Get.new(uri)
  headers.each { |k, v| req[k] = v }

  res = Net::HTTP.start(uri.hostname, uri.port) {|http|
    http.request(req)
  }
  JSON.parse(res.body)
end

def getScores(ids, headers)
  uri = URI("http://scores.base.svc.cluster.local:7000/scores?ids=#{ids.join(',')}")
  req = Net::HTTP::Get.new(uri)
  headers.each { |k, v| req[k] = v }

  res = Net::HTTP.start(uri.hostname, uri.port) {|http|
    http.request(req)
  }
  JSON.parse(res.body)
end

def getForwardHeaders(env)
  headers = {}
  forwardHeaders = [
    "cookie",
    'x-request-id',
    'x-b3-traceid',
    'x-b3-spanid',
    'x-b3-parentspanid',
    'x-b3-sampled',
    'x-b3-flags',
    'x-ot-span-context'
  ]

  forwardHeaders.each do |h|
    envKey = "HTTP_#{h.upcase.gsub('-', '_')}"
    if env[envKey]
      headers[h] = env[envKey]
    end
  end

  return headers
end
