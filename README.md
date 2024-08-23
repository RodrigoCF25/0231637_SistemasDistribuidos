PUERTO: 8080

USAR POSTMAN

Mandar un JSON al segmento write (http://localhost:8080/write)

    SE REQUIERE MÉTODO POST

    El JSON que es válido, ejemplo:

    {
        "value": "TGV0J3MgR28gIzEK"
    }

    Cuando se inserta el Record, se imprime el record en Consola. Todo record que se busca agregar, su offset es modificado para ser el último del slice

Mandar un JSON al segmento read (http://localhost:8080/read)

    SE REQUIERE MÉTODO GET

    EL JSON que es válido, ejemplo:

    {

        "offset": 0
    }

    Imprime en Consola el Record obtenido si el offset existe y un record como JSON por http, ejemplo: {"value:"TGV0J3MgR28gIzEK"}


