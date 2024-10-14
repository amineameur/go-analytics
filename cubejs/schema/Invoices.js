cube(`Invoices`, {
    sql: `SELECT * FROM invoices`,
  
    measures: {
      count: {
        type: `count`,
        drillMembers: [id, createdAt]
      },
      totalAmount: {
        sql: `amount`,
        type: `sum`
      },
      averageAmount: {
        sql: `amount`,
        type: `avg`
      }
    },
  
    dimensions: {
      id: {
        sql: `id`,
        type: `number`,
        primaryKey: true
      },
      createdAt: {
        sql: `created_at`,
        type: `time`
      }
    }
  });
  