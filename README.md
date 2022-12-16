# Article Service
The Article Service is a microservice that provides a HTTP and GRPC interface for creating and find article by id
## HTTP Server
The HTTP server provides a `POST` endpoint at `/articles` for creating new articles. The request body should be a JSON object with the following fields:
- title: the title of the article (string)
- author: the author of the article (string)
- body: the body of the article (string)

Example request:

```
POST /articles
Content-Type: application/json

{
  "title": "My Article",
  "author": "John Doe",
  "body": "This is the body of my article."
}
```

Example response
```
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": 1,
  "title": "My Article",
  "author": "John Doe",
  "body": "This is the body of my article."
}

```

# Search Service
The Search Service is a microservice that provides a HTTP interface for searching articles.

## HTTP Server
The HTTP server provides a `GET` endpoint at `/articles` for searching articles. The request can include the following query parameters:

size: the number of results to return per page (integer, default: 10)
page: the page number of the results (integer, default: 1)
author: the author of the articles to search for (string)

Example request:
```
GET /articles?size=10&page=1&author=Kadhafi2
```

Example response:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "articles": [
    {
      "id": 35,
      "author": "Irvan Kadhafi2",
      "title": "Teh Obeng Batam, Minuman Tradisional dari Kota Industri",
      "body": "Teh obeng adalah minuman tradisional yang berasal dari Kota Batam, Kepulauan Riau. Teh obeng terbuat dari campuran teh hitam, gula merah, dan rempah-rempah seperti kayu manis, cengkih, dan lain-lain. Minuman ini dianggap sebagai salah satu ciri khas Batam dan sering dijual di pasar-pasar tradisional di kota tersebut.\n\nTeh obeng dianggap sebagai minuman yang berguna untuk menghangatkan tubuh dan meningkatkan daya tahan tubuh. Selain itu, teh obeng juga dianggap sebagai minuman yang bermanfaat untuk mengurangi rasa sakit kepala dan meredakan batuk.\n\nDi Kota Batam, teh obeng juga dianggap sebagai minuman yang cocok untuk disajikan saat santai bersama teman-teman atau keluarga. Selain itu, teh obeng juga sering dijadikan sebagai cinderamata untuk diberikan kepada tamu yang datang berkunjung.\n\nTeh obeng merupakan salah satu minuman tradisional yang sangat populer di Kota Batam dan merupakan bagian yang tidak terpisahkan dari budaya masyarakat setempat. Jika berkunjung ke Kota Batam, jangan lupa untuk mencicipi teh obeng sebagai salah satu kenangan yang tak terlupakan dari kota industri tersebut.",
      "created_at": "2022-12-16T17:27:14.140349Z"
    },
    {
      "id": 34,
      "author": "Irvan Kadhafi2",
      "title": "Teh Obeng Batam, Minuman Tradisional dari Kota Industri",
      "body": "Teh obeng adalah minuman tradisional yang berasal dari Kota Batam, Kepulauan Riau. Teh obeng terbuat dari campuran teh hitam, gula merah, dan rempah-rempah seperti kayu manis, cengkih, dan lain-lain. Minuman ini dianggap sebagai salah satu ciri khas Batam dan sering dijual di pasar-pasar tradisional di kota tersebut.\n\nTeh obeng dianggap sebagai minuman yang berguna untuk menghangatkan tubuh dan meningkatkan daya tahan tubuh. Selain itu, teh obeng juga dianggap sebagai minuman yang bermanfaat untuk mengurangi rasa sakit kepala dan meredakan batuk.\n\nDi Kota Batam, teh obeng juga dianggap sebagai minuman yang cocok untuk disajikan saat santai bersama teman-teman atau keluarga. Selain itu, teh obeng juga sering dijadikan sebagai cinderamata untuk diberikan kepada tamu yang datang berkunjung.\n\nTeh obeng merupakan salah satu minuman tradisional yang sangat populer di Kota Batam dan merupakan bagian yang tidak terpisahkan dari budaya masyarakat setempat. Jika berkunjung ke Kota Batam, jangan lupa untuk mencicipi teh obeng sebagai salah satu kenangan yang tak terlupakan dari kota industri tersebut.",
      "created_at": "2022-12-16T17:25:59.61041Z"
    },
    {
      "id": 33,
      "author": "Irvan Kadhafi2",
      "title": "Teh Obeng Batam, Minuman Tradisional dari Kota Industri",
      "body": "Teh obeng adalah minuman tradisional yang berasal dari Kota Batam, Kepulauan Riau. Teh obeng terbuat dari campuran teh hitam, gula merah, dan rempah-rempah seperti kayu manis, cengkih, dan lain-lain. Minuman ini dianggap sebagai salah satu ciri khas Batam dan sering dijual di pasar-pasar tradisional di kota tersebut.\n\nTeh obeng dianggap sebagai minuman yang berguna untuk menghangatkan tubuh dan meningkatkan daya tahan tubuh. Selain itu, teh obeng juga dianggap sebagai minuman yang bermanfaat untuk mengurangi rasa sakit kepala dan meredakan batuk.\n\nDi Kota Batam, teh obeng juga dianggap sebagai minuman yang cocok untuk disajikan saat santai bersama teman-teman atau keluarga. Selain itu, teh obeng juga sering dijadikan sebagai cinderamata untuk diberikan kepada tamu yang datang berkunjung.\n\nTeh obeng merupakan salah satu minuman tradisional yang sangat populer di Kota Batam dan merupakan bagian yang tidak terpisahkan dari budaya masyarakat setempat. Jika berkunjung ke Kota Batam, jangan lupa untuk mencicipi teh obeng sebagai salah satu kenangan yang tak terlupakan dari kota industri tersebut.",
      "created_at": "2022-12-16T14:57:42.26743Z"
    },
    {
      "id": 32,
      "author": "Irvan Kadhafi2",
      "title": "Teh Obeng Batam, Minuman Tradisional dari Kota Industri",
      "body": "Teh obeng adalah minuman tradisional yang berasal dari Kota Batam, Kepulauan Riau. Teh obeng terbuat dari campuran teh hitam, gula merah, dan rempah-rempah seperti kayu manis, cengkih, dan lain-lain. Minuman ini dianggap sebagai salah satu ciri khas Batam dan sering dijual di pasar-pasar tradisional di kota tersebut.\n\nTeh obeng dianggap sebagai minuman yang berguna untuk menghangatkan tubuh dan meningkatkan daya tahan tubuh. Selain itu, teh obeng juga dianggap sebagai minuman yang bermanfaat untuk mengurangi rasa sakit kepala dan meredakan batuk.\n\nDi Kota Batam, teh obeng juga dianggap sebagai minuman yang cocok untuk disajikan saat santai bersama teman-teman atau keluarga. Selain itu, teh obeng juga sering dijadikan sebagai cinderamata untuk diberikan kepada tamu yang datang berkunjung.\n\nTeh obeng merupakan salah satu minuman tradisional yang sangat populer di Kota Batam dan merupakan bagian yang tidak terpisahkan dari budaya masyarakat setempat. Jika berkunjung ke Kota Batam, jangan lupa untuk mencicipi teh obeng sebagai salah satu kenangan yang tak terlupakan dari kota industri tersebut.",
      "created_at": "2022-12-16T14:53:00.002766Z"
    },
    {
      "id": 31,
      "author": "Irvan Kadhafi2",
      "title": "Teh Obeng Batam, Minuman Tradisional dari Kota Industri",
      "body": "Teh obeng adalah minuman tradisional yang berasal dari Kota Batam, Kepulauan Riau. Teh obeng terbuat dari campuran teh hitam, gula merah, dan rempah-rempah seperti kayu manis, cengkih, dan lain-lain. Minuman ini dianggap sebagai salah satu ciri khas Batam dan sering dijual di pasar-pasar tradisional di kota tersebut.\n\nTeh obeng dianggap sebagai minuman yang berguna untuk menghangatkan tubuh dan meningkatkan daya tahan tubuh. Selain itu, teh obeng juga dianggap sebagai minuman yang bermanfaat untuk mengurangi rasa sakit kepala dan meredakan batuk.\n\nDi Kota Batam, teh obeng juga dianggap sebagai minuman yang cocok untuk disajikan saat santai bersama teman-teman atau keluarga. Selain itu, teh obeng juga sering dijadikan sebagai cinderamata untuk diberikan kepada tamu yang datang berkunjung.\n\nTeh obeng merupakan salah satu minuman tradisional yang sangat populer di Kota Batam dan merupakan bagian yang tidak terpisahkan dari budaya masyarakat setempat. Jika berkunjung ke Kota Batam, jangan lupa untuk mencicipi teh obeng sebagai salah satu kenangan yang tak terlupakan dari kota industri tersebut.",
      "created_at": "2022-12-16T14:52:10.001681Z"
    }
  ],
  "cursor_info": {
    "size": 10,
    "count": 5,
    "countPage": 1,
    "hasMore": false,
    "cursor": "1"
  }
}
```

# Gateway Service (TODO)
The Gateway Service is a microservice that acts as a single entry point for all requests to the Article Service and Search Service. It is responsible for routing requests to the appropriate service and consolidating the responses from multiple services into a single response for the client.

# TODO
- Implement unit tests
- Consolidate the endpoints of the Article Service and Search Service into a single endpoint with the Gateway Service