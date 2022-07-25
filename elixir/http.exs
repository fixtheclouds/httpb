Mix.install([{:simplehttp, "~> 0.5.1"}])

[uri] = System.argv

start_at = Time.utc_now
case SimpleHttp.get(uri) do
  {:ok, _ } ->
    elapsed = Time.diff(Time.utc_now, start_at, :millisecond)
    IO.puts("Fetched in #{elapsed} ms")

  {:error, reason } ->
    IO.puts("Failed: #{reason}")
end
