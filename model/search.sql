select flat.id as "id",
       flat."cost",
       flat."space",
       flat.room_count,
       flat.floor,
       flat."number",
       street."name",
       city."name",
       district."name"
from flat,
     house,
     district,
     street,
     city
where flat.house_id = house.id
  and house.city_id = city.id
  and house.district_id = district.id
  and house.street_id = street.id
limit 100