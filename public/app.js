const form = document.getElementById('converter-form');
const amountInput = document.getElementById('amount');
const fromCurrencyInput = document.getElementById('fromCurrency');
const toCurrencyInput = document.getElementById('toCurrency');
const swapButton = document.getElementById('swap-btn');
const statusEl = document.getElementById('status');
const resultCard = document.getElementById('result-card');
const resultValue = document.getElementById('result-value');
const rateValue = document.getElementById('rate-value');
const targetCurrency = document.getElementById('target-currency');

function setStatus(message, isError = false) {
    statusEl.textContent = message;
    statusEl.style.background = isError
        ? 'rgba(248, 113, 113, 0.14)'
        : 'rgba(56, 189, 248, 0.12)';
    statusEl.style.borderColor = isError
        ? 'rgba(248, 113, 113, 0.3)'
        : 'rgba(56, 189, 248, 0.25)';
}

function formatCurrency(value, currency) {
    try {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency,
            maximumFractionDigits: 2,
        }).format(value);
    } catch {
        return `${currency} ${value.toFixed(2)}`;
    }
}

async function loadCurrencies() {
    try {
        const response = await fetch('/currencies');
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Unable to load currencies.');
        }

        const currencies = data.currencies || [];
        if (!currencies.length) {
            throw new Error('No currencies returned by the server.');
        }

        const optionsHtml = currencies.map((currency) => {
            const label = currency.name ? `${currency.code} - ${currency.name}` : currency.code;
            return `<option value="${currency.code}">${label}</option>`;
        }).join('');

        fromCurrencyInput.innerHTML = optionsHtml;
        toCurrencyInput.innerHTML = optionsHtml;
        fromCurrencyInput.value = 'USD';
        toCurrencyInput.value = 'EUR';
        fromCurrencyInput.disabled = false;
        toCurrencyInput.disabled = false;
        setStatus('Enter an amount and pick your currencies.');
    } catch (error) {
        setStatus(error.message, true);
    }
}

document.addEventListener('DOMContentLoaded', loadCurrencies);

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const amount = Number(amountInput.value);
    const fromCurrency = (fromCurrencyInput.value || '').toUpperCase();
    const toCurrency = (toCurrencyInput.value || '').toUpperCase();

    if (!amount || amount <= 0) {
        setStatus('Please enter a valid amount greater than zero.', true);
        return;
    }

    setStatus('Fetching the latest exchange rate...');
    resultCard.classList.add('hidden');

    try {
        const response = await fetch('/convert', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                from_currency: fromCurrency,
                to_currency: toCurrency,
                amount,
            }),
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Conversion failed.');
        }

        resultValue.textContent = formatCurrency(data.total_converted_amount, data.to_currency);
        rateValue.textContent = data.exchange_rate.toFixed(4);
        targetCurrency.textContent = data.to_currency;
        resultCard.classList.remove('hidden');
        setStatus(`Converted ${formatCurrency(amount, fromCurrency)} to ${data.to_currency}.`);
    } catch (error) {
        setStatus(error.message, true);
    }
});

swapButton.addEventListener('click', () => {
    const fromValue = fromCurrencyInput.value;
    const toValue = toCurrencyInput.value;
    fromCurrencyInput.value = toValue;
    toCurrencyInput.value = fromValue;
});