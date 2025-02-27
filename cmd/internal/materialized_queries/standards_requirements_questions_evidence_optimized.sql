WITH
reference_values_data AS (
    SELECT
        id
        , json_object(
            'id', id
            , 'type_id', type_id
            , 'code', code
            , 'name', name
            , 'description', description
            , 'is_active', is_active
            , 'created_at', created_at
            , 'updated_at', updated_at
        ) AS ref_json
    FROM reference_values
)

, evidence_data AS (
    SELECT
        question_id
        , json_arrayagg(
            json_object(
                'id', id
                , 'question_id', question_id
                , 'type', (
                    SELECT ref_json FROM reference_values_data
                    WHERE id = type_id
                )
                , 'expected', expected
                , 'status', (
                    SELECT ref_json FROM reference_values_data
                    WHERE id = status_id
                )
                , 'created_at', created_at
                , 'updated_at', updated_at
            )
        ) AS evidence_json
    FROM evidence
    GROUP BY question_id
)

, questions_data AS (
    SELECT
        requirement_id
        , json_arrayagg(
            json_object(
                'id', id
                , 'requirement_id', requirement_id
                , 'question', question
                , 'guidance', guidance
                , 'created_at', created_at
                , 'updated_at', updated_at
                , 'evidence', (
                    SELECT evidence_json FROM evidence_data
                    WHERE question_id = id
                )
            )
        ) AS questions_json
    FROM questions
    GROUP BY requirement_id
)

, requirements_data AS (
    SELECT
        standard_id
        , json_arrayagg(
            json_object(
                'id', id
                , 'standard_id', standard_id
                , 'level_id', level_id
                , 'parent_id', parent_id
                , 'reference_code', reference_code
                , 'name', name
                , 'description', description
                , 'questions', (
                    SELECT questions_json FROM questions_data
                    WHERE requirement_id = id
                )
            )
        ) AS requirements_json
    FROM requirement
    GROUP BY standard_id
)

SELECT
    json_pretty(
        json_object(
            'id', s.id
            , 'name', s.name
            , 'description', s.description
            , 'version', s.version
            , 'requirements', (
                SELECT requirements_json FROM requirements_data
                WHERE standard_id = s.id
            )
        )
    ) AS `data`
FROM standards AS s
WHERE s.id = 1;
