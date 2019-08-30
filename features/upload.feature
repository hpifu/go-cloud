Feature: upload 测试

    Scenario Outline: upload
        When 请求 /upload, file: "<file>"
        Then 检查状态码 res.status_code: <status>
        Then 检查 data 目录存在文件, file: "<file>"
        Examples:
            | file          | status |
            | hatlonely.png | 200    |