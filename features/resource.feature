Feature: resource 测试

    Scenario Outline: resource
        When 请求 /resource, file: "<file>"
        Examples:
            | file          |
            | hatlonely.png |