async function loadOrder() {
    const orderId = document.getElementById("orderId").value.trim();
    if (!orderId) {
        document.getElementById("order").innerHTML = "<p style='color:red;'>Введите order_uid!</p>";
        return;
    }

    try {
        const response = await fetch(`http://localhost:8080/order/${orderId}`);
        if (!response.ok) {
            throw new Error("Ошибка при загрузке заказа: " + response.status);
        }
        const order = await response.json();
        renderOrder(order);
    } catch (err) {
        document.getElementById("order").innerHTML = `<p style="color:red;">${err.message}</p>`;
    }
}

function renderOrder(order) {
    const container = document.getElementById("order");

    container.innerHTML = `
    <h2>Основная информация</h2>
    <p><b>UID заказа:</b> ${order.order_uid}</p>
    <p><b>Трек-номер:</b> ${order.track_number}</p>
    <p><b>Entry:</b> ${order.entry}</p>
    <p><b>Locale:</b> ${order.locale}</p>
    <p><b>Internal Signature:</b> ${order.internal_signature}</p>
    <p><b>Customer ID:</b> ${order.customer_id}</p>
    <p><b>Delivery Service:</b> ${order.delivery_service}</p>
    <p><b>ShardKey:</b> ${order.shardkey}</p>
    <p><b>SmID:</b> ${order.sm_id}</p>
    <p><b>OofShard:</b> ${order.oof_shard}</p>
    <p><b>Дата создания:</b> ${new Date(order.date_created).toLocaleString()}</p>

    <h2>Доставка</h2>
    <p><b>Имя:</b> ${order.delivery.name}</p>
    <p><b>Телефон:</b> ${order.delivery.phone}</p>
    <p><b>Email:</b> ${order.delivery.email}</p>
    <p><b>Город:</b> ${order.delivery.city}</p>
    <p><b>Адрес:</b> ${order.delivery.address}</p>
    <p><b>Регион:</b> ${order.delivery.region}</p>
    <p><b>ZIP:</b> ${order.delivery.zip}</p>

    <h2>Оплата</h2>
    <p><b>Transaction:</b> ${order.payment.transaction}</p>
    <p><b>Request ID:</b> ${order.payment.request_id}</p>
    <p><b>Сумма:</b> ${order.payment.amount} ${order.payment.currency}</p>
    <p><b>PaymentDT:</b> ${order.payment.payment_dt}</p>
    <p><b>Банк:</b> ${order.payment.bank}</p>
    <p><b>Провайдер:</b> ${order.payment.provider}</p>
    <p><b>Стоимость доставки:</b> ${order.payment.delivery_cost}</p>
    <p><b>Итого за товары:</b> ${order.payment.goods_total}</p>
    <p><b>Custom Fee:</b> ${order.payment.custom_fee}</p>

    <h2>Товары</h2>
    <table class="items-table" border="1" cellpadding="5">
      <tr>
        <th>ChrtID</th>
        <th>TrackNumber</th>
        <th>RID</th>
        <th>Название</th>
        <th>Бренд</th>
        <th>Цена</th>
        <th>Скидка</th>
        <th>Размер</th>
        <th>Итого</th>
        <th>NmID</th>
        <th>Status</th>
      </tr>
      ${order.items.map(item => `
        <tr>
          <td>${item.chrt_id}</td>
          <td>${item.track_number}</td>
          <td>${item.rid}</td>
          <td>${item.name}</td>
          <td>${item.brand}</td>
          <td>${item.price}</td>
          <td>${item.sale}%</td>
          <td>${item.size}</td>
          <td>${item.total_price}</td>
          <td>${item.nm_id}</td>
          <td>${item.status}</td>
        </tr>
      `).join("")}
    </table>
  `;
}

