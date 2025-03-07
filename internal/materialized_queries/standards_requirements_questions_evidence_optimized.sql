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
            , 'deleted_at', deleted_at
        ) AS ref_json
    FROM reference_values
)

, evidence_data AS (
    SELECT
        e.question_id
        , json_arrayagg(
            json_object(
                'id', e.id
                , 'question_id', e.question_id
                , 'type', (
                    SELECT rv.ref_json FROM reference_values_data AS rv
                    WHERE rv.id = e.type_id
                )
                , 'expected', e.expected
                , 'created_at', e.created_at
                , 'updated_at', e.updated_at
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
                'id', q.id
                , 'requirement_id', q.requirement_id
                , 'question', q.question
                , 'guidance', q.guidance
                , 'created_at', q.created_at
                , 'updated_at', q.updated_at
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
                'id', r.id
                , 'standard_id', r.standard_id
                , 'level_id', r.level_id
                , 'parent_id', r.parent_id
                , 'reference_code', r.reference_code
                , 'name', r.name
                , 'description', r.description
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
        'id', s.id
        , 'name', s.name
        , 'description', s.description
        , 'version', s.version
        , 'requirements', (
            SELECT rd.requirements_json FROM requirements_data AS rd
            WHERE rd.standard_id = s.id
        )
    ) AS `data`
    -- ) AS `data`
FROM standards AS s
WHERE s.id = 1;
