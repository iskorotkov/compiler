program p;
var i: integer;
    s: double;
function foo(a: integer): integer;
begin
  foo := a;
end
function square(x: double): double;
begin
  square := x * x;
end
begin
  i := foo(1);
  s := square(2.0);
end.
