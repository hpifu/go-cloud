Feature: GET /resource

    Scenario: case
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123
            }
            """
        When http 请求 POST /upload/d571bda90c2d4e32a793b8a1ff4ff984
            """
            {
                "file": "features/assets/hatlonely.png"
            }
            """
        Then http 检查 200
        When http 请求 GET /resource/d571bda90c2d4e32a793b8a1ff4ff984
            """
            {
                "params": {
                    "name": "hatlonely.png"
                }
            }
            """
        Then http 检查 200
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
