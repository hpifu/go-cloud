Feature: POST /upload/:id

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
        Then fs 检查文件存在 "output/go-cloud/data/123/hatlonely.png"
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case directory
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
                "file": "features/assets/hatlonely.png",
                "params": {
                    "name": "/_account/1.png"
                }
            }
            """
        Then http 检查 200
        Then fs 检查文件存在 "output/go-cloud/data/123/_account/1.png"
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"

    Scenario: case directory
        Given redis set object "d571bda90c2d4e32a793b8a1ff4ff984"
            """
            {
                "id": 123
            }
            """
        When http 请求 POST /upload/124
            """
            {
                "header": {
                    "Authorization": "d571bda90c2d4e32a793b8a1ff4ff984"
                },
                "file": "features/assets/hatlonely.png",
                "params": {
                    "name": "/_account/1.png"
                }
            }
            """
        Then http 检查 403
            """
            {
                "text": "您没有该资源的权限"
            }
            """
        Then fs 检查文件存在 "output/go-cloud/data/123/_account/1.png"
        Given redis del "d571bda90c2d4e32a793b8a1ff4ff984"
