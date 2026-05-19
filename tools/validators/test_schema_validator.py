import os
import sys
import unittest


_CURRENT_DIR = os.path.dirname(__file__)
if _CURRENT_DIR not in sys.path:
    sys.path.insert(0, _CURRENT_DIR)


from schema_validator import SchemaValidator

class SchemaConstraintValidateUnittest(unittest.TestCase):

    def setUp(self):
        base_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), "../.."))
        expanded_schema_dir = os.path.abspath(os.path.join(base_dir, "expanded_schemas"))
        base_schema_path = os.path.abspath(os.path.join(base_dir, "schemas/base.yaml"))
        self._schema_validator = SchemaValidator(expanded_schema_dir, base_schema_path, console_log=True)

    def test_required_constraint(self):
        constraint = {
            "required": True
        }
        validate_result = {
            "errors": []
        }
        self.assertEqual(len(validate_result["errors"]), 0)
        self._schema_validator._validate_constraints("", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints("123", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints({}, constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)
        self._schema_validator._validate_constraints(None, constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 3)

    def test_pattern_constraint(self):
        constraint = {
            "pattern": '\\w*-\\d'
        }
        validate_result = {
            "errors": []
        }
        self._schema_validator._validate_constraints("abc-1", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 0)
        self._schema_validator._validate_constraints("abc-12", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 0)
        self._schema_validator._validate_constraints("$^", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        constraint = {
            "pattern": '^\\w*-\\d$'
        }
        self._schema_validator._validate_constraints("abc-12", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)

    def test_max_length_constraint(self):
        constraint = {
            "max_len": 5
        }
        validate_result = {
            "errors": []
        }
        self._schema_validator._validate_constraints("12345", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 0)
        self._schema_validator._validate_constraints("123456", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)

        constraint = {
            "min_len": 2
        }
        self._schema_validator._validate_constraints("1", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)
        self._schema_validator._validate_constraints("12", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)

        constraint = {
            "min_len": 2,
            "max_len": 5
        }
        self._schema_validator._validate_constraints("1234", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)
        self._schema_validator._validate_constraints("1", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 3)
        self._schema_validator._validate_constraints("1123456", constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 4)

    def test_array_length_constraint(self):
        constraint = {
            "array": {
                "max_size": 3,
                "min_size": 2
            }
        }
        validate_result = {
            "errors": []
        }
        self._schema_validator._validate_constraints([{"a": "b"}, {"c": "d"}], constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 0)
        self._schema_validator._validate_constraints([{"a": "b"}], constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints([{"a": "b"}, {"c": "d"}, {"c": "d"}, {"c": "d"}], constraint, "field_path",
                                                     validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)

    def test_enum_constraint(self):
        constraint = {
            "enum": {
                "values": [
                    'a',
                    'b',
                    'c'
                ]
            }
        }
        validate_result = {
            "errors": []
        }
        self._schema_validator._validate_constraints('a', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 0)
        self._schema_validator._validate_constraints('d', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        constraint = {
            "enum": {
                "values": [
                    'a',
                    'b',
                    'c'
                ],
                "default_value": 'a'
            }
        }
        validate_result = {
            "errors": []
        }
        self._schema_validator._validate_constraints('', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints('a', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints('b', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints('d', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 2)

    def test_oneof_constraint(self):
        constraint = {
            "oneOf": [
                {
                    "type": 'semantic_string',
                },
            ]
        }
        validate_result = {
            "errors": []
        }
        self._schema_validator._validate_constraints(1, constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints('a', constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)
        self._schema_validator._validate_constraints({
            "a": "b"
        }, constraint, "field_path", validate_result)
        self.assertEqual(len(validate_result["errors"]), 1)

if __name__ == '__main__':
    unittest.main()


