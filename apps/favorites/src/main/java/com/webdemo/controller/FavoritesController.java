package com.webdemo.controller;

import com.google.common.collect.Maps;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.math.BigDecimal;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

@RestController
public class FavoritesController {

    @GetMapping("/favorites")
    public Map<Integer, Integer> sales(@RequestParam("ids") String ids) {
        Map<Integer, Integer> mockDbMap = Maps.newHashMap();
        mockDbMap.put(1, 90);
        mockDbMap.put(2, 43);
        mockDbMap.put(3, 44);
        mockDbMap.put(4, 80);
        mockDbMap.put(5, 67);
        mockDbMap.put(6, 39);
        mockDbMap.put(7, 50);
        mockDbMap.put(8, 77);
        mockDbMap.put(9, 89);
        mockDbMap.put(10, 101);
        mockDbMap.put(11, 60);
        mockDbMap.put(12, 51);
        mockDbMap.put(13, 32);
        mockDbMap.put(14, 49);
        mockDbMap.put(15, 9);

        Map<Integer, Integer> resultMap = Maps.newHashMap();

        System.out.println("getting favorites of ids: " + ids);
        for (String id : ids.split(",")) {
            Integer value = mockDbMap.get(Integer.valueOf(id));
            if (value != null) {
                resultMap.put(Integer.valueOf(id), value);
            }
        }

        return resultMap;
    }

}
