-- name: CreateTriage :exec
INSERT INTO triages (id, type, grid, sku_sap, sku_wms, description, cust_id, seller, quantity_supplied, final_quantity, unitary_value, total_value_offered, final_total_value, category, sub_category, sent_to_batch, sent_to_bling, defect, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20);

-- name: GetTriages :many
SELECT id,
        type,
        grid,
        sku_sap,
        sku_wms,
        description,
        cust_id,
        seller,
        quantity_supplied,
        final_quantity,
        unitary_value,
        total_value_offered,
        final_total_value,
        category,
        sub_category,
        sent_to_batch,
        sent_to_bling,
        defect,
        created_at,
        updated_at
FROM triages;

