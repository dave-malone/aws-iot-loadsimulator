# pylint: disable=W0703
"""
create table function
"""
import logging
import pymysql
import cfn_resource

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)


def create_sql_table(table_name, conn):
    """
    create a database table
    :param table_name:
    :param conn:
    :return:
    """
    create_table = 'CREATE TABLE IF NOT EXISTS ' + table_name + '(' \
                   'uuid VARCHAR(36) NOT NULL ,' \
                   'thingid VARCHAR(128) NOT NULL,' \
                   'iot_timestamp BIGINT NOT NULL,' \
                   'intermediate_timestamp BIGINT NOT NULL,' \
                   'db_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,' \
                   'payload JSON NOT NULL,' \
                   'PRIMARY KEY (uuid));'

    with conn.cursor() as cur:
        cur.execute(create_table)
        conn.commit()


# set `handler` as the entry point for Lambda
handler = cfn_resource.Resource()


@handler.create
def create_thing(event, context):
    """
    create method
    :param event:
    :param context:
    :return:
    """
    del context
    rds_host = event['ResourceProperties']['host']
    name = event['ResourceProperties']['username']
    password = event['ResourceProperties']['password']
    db_name = event['ResourceProperties']['dbName']
    sqs_table_name = event['ResourceProperties']['sqsTableName']

    try:
        conn = pymysql.connect(rds_host,
                               user=name,
                               passwd=password,
                               db=db_name,
                               connect_timeout=5,
                               port=3306)
        LOGGER.info("SUCCESS: Connection to RDS MySQL instance succeeded")
    except Exception as ex:
        LOGGER.error("ERROR: Unexpected error: Could not connect to DB... %s",
                     str(ex))
        return {"Status": "SUCCESS", "PhysicalResourceId": "tableCreate"}
    try:
        create_sql_table(sqs_table_name, conn)
        LOGGER.info("Successfully created second sql table")
        return {"Status": "SUCCESS", "PhysicalResourceId": "tableCreate"}
    except Exception as ex:
        LOGGER.error("ERROR: Unexpected error: Could not create table... %s",
                     str(ex))
        return {"Status": "SUCCESS", "PhysicalResourceId": "tableCreate"}


@handler.update
def update_thing(event, context):
    """
    update method
    :param event:
    :param context:
    :return:
    """
    del context
    rds_host = event['ResourceProperties']['host']
    name = event['ResourceProperties']['username']
    password = event['ResourceProperties']['password']
    db_name = event['ResourceProperties']['dbName']
    sqs_table_name = event['ResourceProperties']['sqsTableName']

    try:
        conn = pymysql.connect(rds_host,
                               user=name,
                               passwd=password,
                               db=db_name,
                               connect_timeout=5,
                               port=3306)
        LOGGER.info("SUCCESS: Connection to RDS MySQL instance succeeded")
    except Exception as ex:
        LOGGER.error("ERROR: Unexpected error: Could not connect to DB... %s",
                     str(ex))
        return {"Status": "SUCCESS", "PhysicalResourceId": "tableCreate"}
    try:
        create_sql_table(sqs_table_name, conn)
        LOGGER.info("Successfully created second sql table")
        return {"Status": "SUCCESS", "PhysicalResourceId": "tableCreate"}
    except Exception as ex:
        LOGGER.error("ERROR: Unexpected error: Could not create table... %s",
                     str(ex))
        return {"Status": "SUCCESS", "PhysicalResourceId": "tableCreate"}
