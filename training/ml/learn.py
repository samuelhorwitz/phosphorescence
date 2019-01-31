from keras.models import Sequential
from keras.layers import Dense, BatchNormalization, Dropout
import numpy
import json
import tensorflowjs as tfjs

dataset = numpy.loadtxt("../data/csv/data.csv", delimiter=",", skiprows=1, usecols=(0,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,18,19,21,22,23,24))
numpy.random.shuffle(dataset)
total_rows, total_cols = dataset.shape
train_count = int(total_rows * 0.9)
input_col_count = total_cols - 2
non_averageable_col_count = 13
training_dataset = dataset[:train_count,:]
testing_dataset = dataset[train_count:,:]
averageable_data = training_dataset[:,non_averageable_col_count:input_col_count]
mean = averageable_data.mean(axis=0)
std = averageable_data.std(axis=0)
clean_training_dataset = numpy.concatenate((training_dataset[:,:non_averageable_col_count], (training_dataset[:,non_averageable_col_count:input_col_count] - mean) / std, training_dataset[:,input_col_count:]), axis=1)
clean_testing_dataset = numpy.concatenate((testing_dataset[:,:non_averageable_col_count], (testing_dataset[:,non_averageable_col_count:input_col_count] - mean) / std, testing_dataset[:,input_col_count:]), axis=1)

numpy.savetxt("training_cleaned.csv", clean_training_dataset, delimiter=",")

with open('meanstd.json', 'w') as outfile:
    json.dump({"mean": mean.tolist(), "std": std.tolist()}, outfile)

X = clean_training_dataset[:,0:input_col_count]
Xtest = clean_testing_dataset[:,0:input_col_count]
AY = clean_training_dataset[:,input_col_count:input_col_count+1]
AYtest = clean_testing_dataset[:,input_col_count:input_col_count+1]
PY = clean_training_dataset[:,input_col_count+1]
PYtest = clean_testing_dataset[:,input_col_count+1]

a_model = Sequential()
a_model.add(Dense(12, input_dim=input_col_count, activation='relu'))
a_model.add(BatchNormalization())
a_model.add(Dropout(0.1))
a_model.add(Dense(8, activation='relu'))
a_model.add(BatchNormalization())
a_model.add(Dropout(0.1))
a_model.add(Dense(4, activation='tanh'))
a_model.add(BatchNormalization())
a_model.add(Dropout(0.1))
a_model.add(Dense(1, activation='sigmoid'))

a_model.compile(loss='logcosh', optimizer='adam', metrics=['logcosh'])

a_model.fit(X, AY, epochs=500)

a_scores = a_model.evaluate(Xtest, AYtest)

p_model = Sequential()
p_model.add(Dense(12, input_dim=input_col_count, activation='relu'))
p_model.add(BatchNormalization())
p_model.add(Dropout(0.1))
p_model.add(Dense(8, activation='relu'))
p_model.add(BatchNormalization())
p_model.add(Dropout(0.1))
p_model.add(Dense(4, activation='tanh'))
p_model.add(BatchNormalization())
p_model.add(Dropout(0.1))
p_model.add(Dense(1, activation='sigmoid'))

p_model.compile(loss='logcosh', optimizer='adam', metrics=['logcosh'])

p_model.fit(X, PY, epochs=500)

p_scores = p_model.evaluate(Xtest, PYtest)

print("Aetherealness \n%s: %.2f%%" % (a_model.metrics_names[1], a_scores[1]*100))
print("Primordialness \n%s: %.2f%%" % (p_model.metrics_names[1], p_scores[1]*100))

tfjs.converters.save_keras_model(a_model, './aetherealness')
tfjs.converters.save_keras_model(p_model, './primordialness')
