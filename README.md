# datagen

Datagen is a CLI utility to generate JSON data using fake values.
It generates data to a file, or to Kafka. 

## Usage

### Examples
**Generate data to a file**
```
./datagen create --filesystem --from ./templates/myTemplate.json
```

It generates data using `myTemplate.json` to a file in `./data/data.txt`. 

**Generate data to Kafka**
```
./datagen create --kafka --topic myTopic --from ./templates/myTemplate.json
```

It generates data using `myTemplate.json` to a kafka topic named `myTopic`. The data will continue to be generated until datagen is stopped. 

## Templates
Datagen relies on valid JSON template documents to generate the data. The templates allow for the creation of fake data. Fake data can be specified using the syntax `fake:{<fake specifier>}` where `<fake specifier>` must be replaced with one of:

**Numbers**
```
number:<min>,<max>
int8
int16
int32
int64
uint8
uint16
uint32
uint64
float32
float32range:<min>,<max>
float64
float64range:<min>,<max>
```
**Strings**
```
digit
letter
lexify:<string>
numerify:<string>
```
**Person**
```
person
name
nameprefix
namesuffix
firstname
lastname
gender
ssn
contact
email
phone
phoneformatted
```

...TODO Add more here

### Template Example

```
{
    "id": "fake:{number:1,11}",
    "name": "fake:{uint64}",
    "department": "IT",
    "designation": "Product Manager",
    "username": "fake:{username}",
    "password": "fake:{password}",
    "address1": "fake:{address}",
    "latitude": "fake:{latitude}",
    "longitude": "fake:{longitude}",
    "latitudeInRange": "fake:{latituderange:23.1,56.7}",
    "address": {
        "city": "Mumbai",
        "state": "Maharashtra",
        "country": "India"
    }
}
```

The above template would produce JSON documents like the following:

```
{
    "address": {
        "city": "Mumbai",
        "country": "India",
        "state": "Maharashtra"
    },
    "address1": {
        "address": "New Mexico Stehrtown Valleyschester743号 ",
        "street": "Valleyschester743号",
        "city": "Stehrtown",
        "state": "New Mexico",
        "zip": "36339",
        "country": "New Caledonia",
        "latitude": -49.659845,
        "longitude": 18.110164
    },
    "department": "IT",
    "designation": "Product Manager",
    "id": 5,
    "latitude": -1.947745,
    "latitudeInRange": 29.420478,
    "longitude": -27.482948,
    "name": 5427729115751982900,
    "password": "vGdku\u0026$ZGi8#",
    "username": "Kiehn4009"
}
```
