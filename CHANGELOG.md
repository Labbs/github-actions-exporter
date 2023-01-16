# ØKP4 github-actions-exporter changelog

## 1.0.0 (2023-01-16)


### Features

* add option to disable fetching workflow usage ([03a380a](https://github.com/okp4/github-actions-exporter/commit/03a380a33ee224507e36d23379527577862c0538))
* default to fetching metrics for all repositories ([04e6c31](https://github.com/okp4/github-actions-exporter/commit/04e6c3119258727227306ce58c5b61d9674c247a))
* **newMetrics:** Added bill metrics ([e7ff476](https://github.com/okp4/github-actions-exporter/commit/e7ff476b784522f05cb4d690c20b3f4545e21961))
* paginate everything, handle ratelimit ([e638be9](https://github.com/okp4/github-actions-exporter/commit/e638be99ef4ee5e4fa617f5bf33b6e63da6adc44))


### Bug Fixes

* **authentication:** the status code returned now is checked ([f183899](https://github.com/okp4/github-actions-exporter/commit/f18389997e2f4e18fdcf5d48c31e094b4c156dac))


### Performance Improvements

* fetch only 12hr of workflow runs ([edc20aa](https://github.com/okp4/github-actions-exporter/commit/edc20aaa5a5a49d55b6453c0a464f078f518095f))
* limit cache size ([627d88f](https://github.com/okp4/github-actions-exporter/commit/627d88f769cb1648f40f1509a06bc5a1d5062c65))
* only measure non-empty repositories ([19a8546](https://github.com/okp4/github-actions-exporter/commit/19a85464d047216c584ec976715759abbb61057a))
* use httpcache ([1fe326d](https://github.com/okp4/github-actions-exporter/commit/1fe326d88c1487a5d051d2a8e994fa49851c7ee5))
