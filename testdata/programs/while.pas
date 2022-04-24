program p;
var i, count: integer;
begin
  count := 1;
  i := 0;
  while i < 5 do
  begin
    writeln(i);
    writeln(count);

    i := i + 1;
    count := count * 2;
  end;
  writeln(i);
  writeln(count);
end.
