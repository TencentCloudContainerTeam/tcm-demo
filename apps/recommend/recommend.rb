require 'sinatra'
require 'net/http'

set :bind, '0.0.0.0'
set :port, 7000

get '/recommend' do
  content_type :json
  user = request.env["HTTP_USER"]
  puts "getting recommend products for user #{user}"

  headers = getForwardHeaders(request.env)

  ids = getRecommendIDsByUser(user)
  products = getProducts(ids, headers)
  products.each do |p|
    score = getScores([p["id"]], headers)
    p["score"] = score[p["id"].to_s]
  end
  recommend = {
    products: products,
    #banner: "v2 banner todo" #v2
  }

  recommend.to_json
end

def getRecommendIDsByUser(user)
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
    "user",
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
