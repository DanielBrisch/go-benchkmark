DROP TABLE IF EXISTS resumo_categoria;

CREATE TABLE resumo_categoria AS
SELECT 
    category, 
    SUM(price * quantity) AS total_vendas, 
    COUNT(*) AS num_vendas
FROM vendas
GROUP BY category;

CREATE TABLE top_produtos AS
SELECT 
    product_id, 
    SUM(quantity) AS total_quantidade
FROM vendas
GROUP BY product_id
ORDER BY total_quantidade DESC
LIMIT 10;

CREATE TABLE vendas_mensais AS
SELECT 
    strftime(date, '%Y-%m') AS mes,
    SUM(price * quantity) AS total_mes
FROM vendas
GROUP BY mes
ORDER BY mes;
