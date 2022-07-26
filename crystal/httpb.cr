require "option_parser"

uri = nil
option_parser = OptionParser.parse do |parser|
  parser.on "-u URL", "--url=URL", "URL to fetch" do |str|
    uri = str
  end
end

struct Result
  property status, time

  def initialize(@status : String, @time : Time)
  end
end
results = [] of Result

Signal::INT.trap do
  print_results(results)
  exit 130
end

def print_results(results)
  if results.size == 0
    puts "Nothing to print..."
    return
  end
  durations = results.map(&.time)

  puts "\nBechmarking HTTP GET to #{uri}"
  puts "avg #{durations.sum(0.0) / durations.size}"
  puts "min #{durations.min}"
  puts "max #{durations.max}"
  puts "#{results.count} requests total"
end

while true
  time = Time.now
  response = Net::HTTP.get_response(URL)
  response.body
  diff = Time.now - time
  puts "fetched in #{diff} seconds, code: #{response.code}"
  results << [diff, response.code]
end
