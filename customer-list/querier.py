from psycopg import Connection


def querier(conn: Connection) -> list[str]:
    customer_list: list[str] = []
    with conn.cursor() as cursor:
        cursor.execute("SELECT DISTINCT customer_number FROM transactions")
        for record in cursor:
            customer_list.append(*record)

        conn.commit()

    return customer_list
