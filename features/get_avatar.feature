Feature: GET /avatar/:id

    Scenario: case
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123
            }
            """
        When http 请求 POST /avatar/123
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "file": "features/assets/hatlonely.png",
                "params": {
                    "name": "1.png"
                }
            }
            """
        Then http 检查 200
        Then fs 检查文件存在 "output/go-cloud/data/123/_pub/account/avatar/1.png"
        When http 请求 GET /avatar/123
            """
            {
                "params": {
                    "name": "1.png"
                }
            }
            """
        Then http 检查 200
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"