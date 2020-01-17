# pylint: disable=W0703
"""
lambda function
"""
import datetime
import json
import logging
import os
import sys
import pymysql

RDS_HOST = os.environ['HOST']
NAME = os.environ['USERNAME']
PASSWORD = os.environ['PASS']
DB_NAME = os.environ['DB_NAME']
TABLE_NAME = os.environ['TABLE_NAME_SQS']

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)

try:
    CONN = pymysql.connect(RDS_HOST,
                           user=NAME,
                           passwd=PASSWORD,
                           db=DB_NAME,
                           connect_timeout=5,
                           port=3306)
except Exception as ex:
    LOGGER.error("ERROR: Unexpected error: Could not connect to DB... %s",
                 str(ex))
    sys.exit()

LOGGER.info("SUCCESS: Connection to RDS MySQL instance succeeded")


def lambda_handler(event, context):
    """
    lambda handler
    :param event:
    :param context:
    :return:
    """
    del context
    LOGGER.info('EVENT: %s', json.dumps(event))
    items = []
    for record in event['Records']:
        # SQS Messages can be read directly
        payload = json.loads(record['body'])
        payload['timestamp'] = datetime.datetime.utcfromtimestamp(
            payload.get('timestamp') / 1000).strftime('%Y-%m-%d %H:%M:%S.%f'),
        payload['messageTs'] = datetime.datetime.utcnow().strftime(
            '%Y-%m-%d %H:%M:%S.%f')
        payload['eventType'] = 'message-success'

        LOGGER.info("Decoded payload: %s", str(payload))
        items.append(payload)

    insert = 'INSERT INTO ' + \
             TABLE_NAME + \
             '(uuid, thingid, iot_timestamp, intermediate_timestamp, payload)' + \
             ' values(%s,%s,%s,%s,%s)'

    with CONN.cursor() as cur:
        for entry in items:
            try:
                cur.execute(
                    insert,
                    (entry['uuid'], entry['thingid'], entry['timestamp'],
                     entry['messageTs'], str(json.dumps(entry['payload']))))
                CONN.commit()
            except Exception as ex:
                print(
                    "ERROR: Unexpected error: Could not add row to table: %s",
                    str(ex))
    CONN.commit()

    return {'status_code': 200, 'body': 'SUCCESS'}
