# Masks
id_length = 63
timestamp_length = 41
node_id_length = 10
sequence_length = 12

timestamp_mask = ('1' * timestamp_length) + ('0' * (id_length - timestamp_length))
node_id_mask = ('0' * (timestamp_length)) + ('1' * node_id_length) + ('0' * (id_length - timestamp_length - node_id_length))
sequence_mask = ('0' * (timestamp_length + node_id_length)) + ('1' * sequence_length) + ('0' * (id_length - timestamp_length - node_id_length - sequence_length))

print(f"Binary Mask for Timestamp should be: {timestamp_mask}")
print(f"Binary Mask for NodeID should be   : {node_id_mask}")
print(f"Binary Mask for Sequence should be : {sequence_mask}")

# Shift Constants
