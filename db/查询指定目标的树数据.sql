SELECT item.*,-CAST(closure.depth AS INTEGER) AS depth
FROM technology_item item
INNER JOIN technology_closure closure
ON item.name = closure.ancestor_name
WHERE closure.descendant_name = "Weaver Organo-Tech Amplifier"

UNION

SELECT item.*,closure.depth
FROM technology_item item
INNER JOIN technology_closure closure
ON item.name = closure.descendant_name 
WHERE closure.ancestor_name = "Weaver Organo-Tech Amplifier"
ORDER BY depth ASC