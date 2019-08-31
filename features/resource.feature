Feature: resource 测试

    Scenario Outline: resource
        When 请求 /upload, file: "<file>"
        When 请求 /resource, file: "<file>"
        Then 检查状态码 res.status_code: <status>
        Examples:
            | file          | status |
            | hatlonely.png | 200    |