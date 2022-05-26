--[[
example useless code to show lua syntax highlighting
this is multiline comment
]]

function blahblahblah(x)

  local table = {
    "asd" = 123,
    "x" = 0.34,  
  }
  if x ~= 3 then
    print( x )
  elseif x == "string"
    my_custom_function( 0x34 )
  else
    unknown_function( "some string" )
  end

  --single line comment
  
end

function blablabla3()

  for k,v in ipairs( table ) do
    --abcde..
    y=[=[
  x=[[
      x is a multi line string
   ]]
  but its definition is iside a highest level string!
  ]=]
    print(" \"\" ")
  --this marks a parser error:
  s = [== asdasdasd]]

    s = math.sin( x )
  end

end

function Fish()
  --[[
end
  --]]
  
  y=[=[
end
  if
  bla
  then
  ble
end
  ]=]
end

function Chips()
  This_line_gets_poorly_indented()
end
