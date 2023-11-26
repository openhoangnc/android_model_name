# Android Model Marketing Name

This repository contains all of the Android model marketing names. The list is based on [this](https://storage.googleapis.com/play_public/supported_devices.html) Google Play Store list.

## Usage

### Kotlin

```kotlin
fun getModelName(manufacturer: String, model: String): String {
    val repoUrl = "https://raw.githubusercontent.com/openhoangnc/android_model_name/main"

    val queryUrl = "$repoUrl/${manufacturer[0]}/$manufacturer/${model.replace("/", "-")}".toLowerCase()

    val url = java.net.URL(queryUrl)

    try {
        url.readText().let {
            return it
        }
    } catch (e: Exception) {
        return model
    }
}

fun getSelfModelName(): String {
    return getModelName(android.os.Build.MANUFACTURER, android.os.Build.MODEL)
}
```
