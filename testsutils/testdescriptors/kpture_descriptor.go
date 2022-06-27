package testdescriptors_test

const ValidKptureDescriptor = `{
    "agents": [
        {
            "metadata": {
                "name": "nginx-deployment-794f656f8b-vtz9s",
                "namespace": "testns2",
                "system": "kubernetes",
                "targetUrl": "10.42.0.204:10000"
            },
            "status": "up",
            "packetNb": 6
        },
        {
            "metadata": {
                "name": "nginx-deployment-794f656f8b-j4ct4",
                "namespace": "testns2",
                "system": "kubernetes",
                "targetUrl": "10.42.0.203:10000"
            },
            "status": "up",
            "packetNb": 6
        },
        {
            "metadata": {
                "name": "nginx-deployment-794f656f8b-xm2ff",
                "namespace": "testns2",
                "system": "kubernetes",
                "targetUrl": "10.42.0.215:10000"
            },
            "status": "up",
            "packetNb": 6
        }
    ],
    "name": "CaptureDemo",
    "uuid": "18388a08-1845-4cfc-891b-29095027babe",
    "captureInfo": {
        "size": 756,
        "packetNb": 18
    },
    "status": "terminated"
}`

const InvalidDescriptor = `{
    "name": "CaptureDemo",
    "uuid": "18388a08-1845-4cfc-891b-29095027babe",
    "captureInfo": {
        "size": 756,
        "packetNb": 18
    "status": "terminated"
}`
