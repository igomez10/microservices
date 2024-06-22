-- delete all tokens that are valid until between 2024-01-22 and 2024-02-22
delete
    FROM public.tokens
    WHERE '2024-01-22' < valid_until AND 
    valid_until < '2024-02-22';


-- vacuum the tokens table to free up space
VACUUM FULL public.tokens;
