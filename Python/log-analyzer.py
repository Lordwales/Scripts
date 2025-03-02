from clickhouse_driver import Client
from kafka import KafkaConsumer
import json

# ClickHouse Configuration
CLICKHOUSE_HOST = "your_clickhouse_host"
DATABASE = "logs"
TABLE = "app_logs"

# Kafka Configuration
KAFKA_BROKER = "your_kafka_broker"
KAFKA_TOPIC = "logs_topic"

# Connect to ClickHouse
client = Client(host=CLICKHOUSE_HOST)

def create_table():
    query = f'''
    CREATE TABLE IF NOT EXISTS {DATABASE}.{TABLE} (
        timestamp DateTime DEFAULT now(),
        level String,
        message String,
        source String
    ) ENGINE = MergeTree()
    ORDER BY timestamp;
    '''
    client.execute(query)
    print("Table ensured in ClickHouse.")

# Function to insert logs into ClickHouse
def insert_log(log):
    query = f"INSERT INTO {DATABASE}.{TABLE} (timestamp, level, message, source) VALUES"
    values = [(log['timestamp'], log['level'], log['message'], log['source'])]
    client.execute(query, values)
    print("Inserted log into ClickHouse.")

# Kafka Consumer Setup
def consume_logs():
    consumer = KafkaConsumer(
        KAFKA_TOPIC,
        bootstrap_servers=KAFKA_BROKER,
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )
    
    print("Consuming logs from Kafka...")
    for message in consumer:
        insert_log(message.value)

if __name__ == "__main__":
    create_table()
    consume_logs()
