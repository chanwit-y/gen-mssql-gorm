-- SELECT
--       s.name AS SchemaName
--      ,t.name AS TableName
--      ,c.name AS ColumnName
--  FROM sys.schemas AS s
--  JOIN sys.tables AS t ON t.schema_id = s.schema_id
--  JOIN sys.columns AS c ON c.object_id = t.object_id
--  ORDER BY
--       SchemaName
--      ,TableName
--      ,ColumnName;
    
 SELECT 
      TABLE_SCHEMA AS SchemaName
     ,TABLE_NAME AS TableName
     ,COLUMN_NAME AS ColumnName
     ,ORDINAL_POSITION AS Position
     ,DATA_TYPE AS DataType
     ,CHARACTER_MAXIMUM_LENGTH AS MaxLeangth
     ,IS_NULLABLE AS IsNullable
 FROM INFORMATION_SCHEMA.COLUMNS
 WHERE TABLE_NAME = 'TRIP_ITEMS'
 ORDER BY
      SchemaName
     ,TableName
     ,ColumnName;



 SELECT TABLE_NAME AS TableName
 FROM INFORMATION_SCHEMA.COLUMNS
 GROUP BY TABLE_NAME
 ORDER BY TableName

 --PKCOLUMN_NAME FKTABLE_NAME

EXEC sp_fkeys @pktable_name = 'TRIP_ITEMS', @pktable_owner = 'dbo'

SELECT CONSTRAINT_NAME AS ConstraintName
       ,CONSTRAINT_TYPE AS Constrainttype 
FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS where TABLE_NAME = 'TRIP_ITEMS'

SELECT 
     KU.table_name as TABLENAME
    ,column_name as PRIMARYKEYCOLUMN
FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC 
INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS KU
    ON TC.CONSTRAINT_TYPE = 'PRIMARY KEY' 
    AND TC.CONSTRAINT_NAME = KU.CONSTRAINT_NAME 
    AND KU.table_name='TRIP_ITEMS'
ORDER BY 
     KU.TABLE_NAME
    ,KU.ORDINAL_POSITION; 