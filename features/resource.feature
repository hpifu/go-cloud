Feature: GET /resource

    Scenario: case
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123
            }
            """
        When http 请求 POST /upload/123
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "file": "features/assets/hatlonely.png"
            }
            """
        Then http 检查 200
        When http 请求 GET /resource/123
            """
            {
                "params": {
                    "name": "hatlonely.png",
                    "token": "d571bda90c2d4e32a793b8a1ff4ff984"
                }
            }
            """
        Then http 检查 200
        When http 请求 GET /resource/123
            """
            {
                "params": {
                    "name": "hatlonely.png"
                },
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                }
            }
            """
        Then http 检查 200
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
