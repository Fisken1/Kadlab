::ignor this 
::@set basePort=8082
::@set /A num = basePort
::echo %basePort%

::@echo docker run -it -p %basePort%:1550 kadlab
::docker run -it -p 8082:1550 kadlab
timeout /t 3 /nobreak
for /l %%x in (1, 1, 5) do (
	timeout /t 3 /nobreak
	::start cmd /k echo test11 + l
   
)

