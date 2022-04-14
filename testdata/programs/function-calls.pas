program p;
var i: integer;
    s: real;
function foo(a: integer): integer;
begin
  foo := a;
end
function square(x: real): real;
begin
  square := x * x;
end
begin
  i := foo(1);
  s := square(2.0);
end.
