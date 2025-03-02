import requests
import jinja2
import pdfkit

# Airtable API configuration
AIRTABLE_API_KEY = "your_airtable_api_key"
BASE_ID = "your_base_id"
TABLE_NAME = "Invoices"
HEADERS = {"Authorization": f"Bearer {AIRTABLE_API_KEY}"}

# Fetch data from Airtable
def fetch_airtable_data():
    url = f"https://api.airtable.com/v0/{BASE_ID}/{TABLE_NAME}"
    response = requests.get(url, headers=HEADERS)
    if response.status_code == 200:
        return response.json()["records"]
    else:
        print("Error fetching data:", response.text)
        return []

# Jinja2 template for invoice
INVOICE_TEMPLATE = """
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .invoice-box { border: 1px solid #ccc; padding: 20px; }
        .header { font-size: 20px; font-weight: bold; }
        .details { margin-top: 20px; }
    </style>
</head>
<body>
    <div class="invoice-box">
        <div class="header">Invoice #{{ invoice['id'] }}</div>
        <p><strong>Client:</strong> {{ invoice['fields']['Client'] }}</p>
        <p><strong>Amount:</strong> ${{ invoice['fields']['Amount'] }}</p>
        <p><strong>Date:</strong> {{ invoice['fields']['Date'] }}</p>
    </div>
</body>
</html>
"""

def generate_invoice_pdf(invoice):
    template = jinja2.Template(INVOICE_TEMPLATE)
    html_content = template.render(invoice=invoice)
    pdf_filename = f"invoice_{invoice['id']}.pdf"
    pdfkit.from_string(html_content, pdf_filename)
    print(f"Generated {pdf_filename}")

# Main execution
def main():
    invoices = fetch_airtable_data()
    for invoice in invoices:
        generate_invoice_pdf(invoice)

if __name__ == "__main__":
    main()
