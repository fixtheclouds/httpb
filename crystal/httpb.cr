require "http/client"
require "option_parser"

uri = ""
failures_count = 0
option_parser = OptionParser.parse do |parser|
  parser.on "-u URL", "--url=URL", "URL to fetch" do |str|
    uri = str
  end
end

struct Result
  property status, time

  def initialize(@status : Int32, @time : Time::Span)
  end
end
results = [] of Result

Signal::INT.trap do
  print_results(results, uri)
  exit 130
end

def print_results(results : Array(Result), uri : String)
  if results.size == 0
    puts "Nothing to print..."
    return
  end
  durations = results.map(&.time).map(&.milliseconds)

  puts "\nBechmarking HTTP GET to #{uri}"
  puts "avg #{durations.sum / durations.size} ms"
  puts "min #{durations.min} ms"
  puts "max #{durations.max} ms"
  puts "#{results.size} requests total"
end

while true
  begin
    time = Time.utc
    response = HTTP::Client.get(uri)
    failures_count += 1
    response.body
    code = response.status_code
    diff = Time.utc - time
    puts "fetched in #{diff.milliseconds} ms, code: #{code}"
    results << Result.new(code, diff)
  rescue Socket::Error
    puts "failed to fetch #{uri}, retrying..."
    failures_count += 1
    sleep 1
  end
end
