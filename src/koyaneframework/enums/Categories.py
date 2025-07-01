from enum import Enum

class StatusCategories(str,Enum):
    ERRORS = "errors"
    STATUS_ANALYZER = "status_analyzer"
    STATUS_GENERATOR = "status_generator"
    SUCCESS_GENERATOR = "success_generator"
    STATUS_BEFORE = "status_before"

class HelpCategories(str,Enum):
    BEFORE = "before"
    GENERATE = "generate"
    EDIT = "edit"
    ANALYZE = "analyze"
