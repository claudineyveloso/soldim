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

-- name: GetTriage :one
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
FROM triages
WHERE id = $1;

UPDATE triages SET type = $2,
        grid = $3,
        sku_sap = $4,
        sku_wms = $5,
        description = $6,
        cust_id = $7,
        seller = $8,
        quantity_supplied = $9,
        final_quantity = $10,
        unitary_value = $11,
        total_value_offered = $12,
        final_total_value = $13,
        category = $14,
        sub_category = $15,
        sent_to_batch = $16,
        sent_to_bling = $17,
        defect = $18,
        updated_at = $19
WHERE triages.id = $1;

