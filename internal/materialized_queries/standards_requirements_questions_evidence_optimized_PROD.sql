WITH
reference_values_data AS (
    SELECT
        id
        , json_object(
            'name', name
        ) AS ref_json
    FROM reference_values
)

, evidence_data AS (
    SELECT
        e.question_id
        , json_arrayagg(
            json_object(
                'type', (
                    SELECT rv.ref_json FROM reference_values_data AS rv
                    WHERE rv.id = e.type_id
                )
                , 'expected', e.expected
            )
        ) AS evidence_json
    FROM evidence AS e
    GROUP BY e.question_id
)

, questions_data AS (
    SELECT
        q.requirement_id
        , json_arrayagg(
            json_object(
                'question', q.question
                , 'evidence', (
                    SELECT ed.evidence_json FROM evidence_data AS ed
                    WHERE ed.question_id = q.id
                )
            )
        ) AS questions_json
    FROM questions AS q
    GROUP BY q.requirement_id
)

, requirements_data AS (
    SELECT
        r.standard_id
        , json_arrayagg(
            json_object(
                'reference_code', r.reference_code
                , 'name', r.name
                , 'questions', (
                    SELECT qd.questions_json FROM questions_data AS qd
                    WHERE qd.requirement_id = r.id
                )
            )
        ) AS requirements_json
    FROM requirement AS r
    GROUP BY r.standard_id
)

SELECT
    -- json_pretty(
    json_object(
        'name', s.name
        , 'version', s.version
        , 'description', s.description
        , 'requirements', (
            SELECT rd.requirements_json FROM requirements_data AS rd
            WHERE rd.standard_id = s.id
        )
    ) AS `data`
    -- ) AS `data`
FROM standards AS s
WHERE s.id = 1;
