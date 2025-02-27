SELECT json_pretty(
    json_object(
        'id', standards.id
        , 'name', standards.name
        , 'description', standards.description
        , 'version', standards.version
        , 'requirements', (
            SELECT
                json_arrayagg(
                    json_object(
                        'id', requirement.id
                        , 'standard_id', requirement.standard_id
                        , 'level_id', requirement.level_id
                        , 'parent_id', requirement.parent_id
                        , 'reference_code', requirement.reference_code
                        , 'name', requirement.name
                        , 'description', requirement.description
                        , 'questions', (
                            SELECT
                                json_arrayagg(
                                    json_object(
                                        'id', questions.id
                                        , 'requirement_id', questions.requirement_id
                                        , 'question', questions.question
                                        , 'guidance', questions.guidance
                                        , 'created_at', questions.created_at
                                        , 'updated_at', questions.updated_at
                                        , 'evidence', (
                                            SELECT
                                                json_arrayagg(
                                                    json_object(
                                                        'id', evidence.id
                                                        , 'question_id', evidence.question_id
                                                        , 'type', (
                                                            SELECT
                                                                json_arrayagg(
                                                                    json_object(
                                                                        'id', reference_values.id
                                                                        , 'type_id', reference_values.type_id
                                                                        , 'code', reference_values.code
                                                                        , 'name', reference_values.name
                                                                        , 'description', reference_values.description
                                                                        , 'is_active', reference_values.is_active
                                                                        , 'created_at', reference_values.created_at
                                                                        , 'updated_at', reference_values.updated_at
                                                                    )
                                                                )
                                                            FROM reference_values
                                                            WHERE reference_values.id = evidence.type_id
                                                        )
                                                        , 'expected', evidence.expected
                                                        , 'status', (
                                                            SELECT
                                                                json_arrayagg(
                                                                    json_object(
                                                                        'id', reference_values.id
                                                                        , 'type_id', reference_values.type_id
                                                                        , 'code', reference_values.code
                                                                        , 'name', reference_values.name
                                                                        , 'description', reference_values.description
                                                                        , 'is_active', reference_values.is_active
                                                                        , 'created_at', reference_values.created_at
                                                                        , 'updated_at', reference_values.updated_at
                                                                    )
                                                                )
                                                            FROM reference_values
                                                            WHERE reference_values.id = evidence.status_id
                                                        )
                                                        , 'created_at', evidence.created_at
                                                        , 'updated_at', evidence.updated_at
                                                    )
                                                )
                                            FROM evidence
                                            WHERE evidence.question_id = questions.id
                                        )
                                    )
                                )
                            FROM questions
                            WHERE questions.requirement_id = requirement.id
                        )
                    )
                )
            FROM requirement
            WHERE requirement.standard_id = standards.id
        )
    )
) AS `data` FROM standards
WHERE standards.id = 1;
