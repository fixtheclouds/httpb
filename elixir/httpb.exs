Mix.install([{:simplehttp, "~> 0.5.1"}])

[uri, number] = System.argv
{limit,_} = Integer.parse(number)

fetch_uri = fn uri ->  SimpleHttp.get(uri) end
do_request = fn _ ->
  start_at = Time.utc_now
  case fetch_uri.(uri) do
    {:ok, _ } ->
      elapsed = Time.diff(Time.utc_now, start_at, :millisecond)
      IO.puts("Fetched in #{elapsed} ms")
      {:ok, elapsed}

    {:error, reason } ->
      IO.puts("Failed: #{reason}")
      {:error, nil}
  end
end

Enum.map(1..limit, do_request)
