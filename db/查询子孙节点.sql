SELECT item.*,closure.depth
FROM technology_item item
INNER JOIN technology_closure closure
ON item.name = closure.descendant_name 
WHERE closure.ancestor_name = "Weaver Organo-Tech Amplifier"
ORDER BY closure.depth ASC