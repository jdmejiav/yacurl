# yacurl en GO

### Instalación GO
Para usarlo, primero de debe instalar Go, el instalador se puede encontrar en la página oficial <a href="https://go.dev/dl/">Instalador</a>

### Modo de uso

Para usalro, primero se debe compilar el archivo yacurl.go

    $ go build yacurl.go

Una vez compilado, se ejecuta pasando como parámetros el hostname y el puerto al que se quiere acceder, de la siguiente manera:

    $ ./yacurl www.google.com 80
