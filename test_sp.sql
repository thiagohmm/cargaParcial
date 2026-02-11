-- Script para verificar se a Stored Procedure existe no banco Oracle

-- 1. Verificar se a SP existe
select object_name,
       object_type,
       status,
       created,
       last_ddl_time
  from user_objects
 where object_name = 'SP_GRAVARINTEGRACAOPRODUTOSTAGING'
   and object_type = 'PROCEDURE';

-- 2. Ver os parâmetros da SP (se existir)
select argument_name,
       data_type,
       in_out,
       position
  from user_arguments
 where object_name = 'SP_GRAVARINTEGRACAOPRODUTOSTAGING'
 order by position;

-- 3. Ver o código da SP (se existir)
select text
  from user_source
 where name = 'SP_GRAVARINTEGRACAOPRODUTOSTAGING'
 order by line;

-- 4. Verificar se a tabela ProdutoIntegracaoStaging existe
select table_name,
       tablespace_name,
       num_rows
  from user_tables
 where table_name = 'PRODUTOINTEGRACAOSTAGING';

-- 5. Ver a estrutura da tabela (se existir)
select column_name,
       data_type,
       data_length,
       nullable,
       data_default
  from user_tab_columns
 where table_name = 'PRODUTOINTEGRACAOSTAGING'
 order by column_id;

-- 6. Teste manual da SP (CUIDADO: vai inserir dados!)
-- BEGIN
--     SP_GRAVARINTEGRACAOPRODUTOSTAGING(p_idRevendedor => 1, p_idProduto => 1);
-- END;
-- /